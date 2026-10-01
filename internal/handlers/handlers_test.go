package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestHealthz(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	h.Healthz(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "{\"status\":\"ok\"}\n" {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func newTestHandler(t *testing.T) (*Handler, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })
	return &Handler{RDB: rdb}, mr
}

func TestCacheSetGet(t *testing.T) {
	h, _ := newTestHandler(t)

	// SET
	req := httptest.NewRequest(http.MethodGet, "/cache/set?key=hello&value=world", nil)
	w := httptest.NewRecorder()
	h.CacheSet(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("set: status = %d, want 200, body=%s", w.Code, w.Body.String())
	}

	// GET
	req = httptest.NewRequest(http.MethodGet, "/cache/get?key=hello", nil)
	w = httptest.NewRecorder()
	h.CacheGet(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get: status = %d, want 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"value":"world"`) {
		t.Fatalf("get: unexpected body: %s", w.Body.String())
	}
}

func TestCacheSetValidation(t *testing.T) {
	h, _ := newTestHandler(t)

	cases := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{"missing key", "/cache/set?value=world", http.StatusBadRequest},
		{"blank key", "/cache/set?key=%20&value=world", http.StatusBadRequest},
		{"missing value", "/cache/set?key=hello", http.StatusBadRequest},
		{"blank value", "/cache/set?key=hello&value=%20", http.StatusBadRequest},
		{"key too long", "/cache/set?key=" + strings.Repeat("a", maxKeyLen+1) + "&value=x", http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			w := httptest.NewRecorder()
			h.CacheSet(w, req)
			if w.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d, body=%s", w.Code, tc.wantStatus, w.Body.String())
			}
		})
	}
}

func TestCacheGetValidation(t *testing.T) {
	h, _ := newTestHandler(t)

	cases := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{"missing key", "/cache/get", http.StatusBadRequest},
		{"blank key", "/cache/get?key=%20", http.StatusBadRequest},
		{"key too long", "/cache/get?key=" + strings.Repeat("a", maxKeyLen+1), http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			w := httptest.NewRecorder()
			h.CacheGet(w, req)
			if w.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", w.Code, tc.wantStatus)
			}
		})
	}
}

func TestCacheGetNotFound(t *testing.T) {
	h, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/cache/get?key=missing", nil)
	w := httptest.NewRecorder()
	h.CacheGet(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}