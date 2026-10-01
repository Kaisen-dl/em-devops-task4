package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateCacheInput(t *testing.T) {
	cases := []struct {
		name         string
		query        string
		requireValue bool
		wantOK       bool
		wantStatus   int
		wantKey      string
		wantValue    string
		wantBodyIn   string
	}{
		{
			name:         "valid key and value",
			query:        "key=hello&value=world",
			requireValue: true,
			wantOK:       true,
			wantKey:      "hello",
			wantValue:    "world",
		},
		{
			name:         "key is trimmed",
			query:        "key=%20%20hello%20%20&value=world",
			requireValue: true,
			wantOK:       true,
			wantKey:      "hello",
			wantValue:    "world",
		},
		{
			name:         "missing key",
			query:        "value=world",
			requireValue: true,
			wantOK:       false,
			wantStatus:   http.StatusBadRequest,
			wantBodyIn:   "key is required",
		},
		{
			name:         "blank key",
			query:        "key=%20%20&value=world",
			requireValue: true,
			wantOK:       false,
			wantStatus:   http.StatusBadRequest,
			wantBodyIn:   "key is required",
		},
		{
			name:         "missing value when required",
			query:        "key=hello",
			requireValue: true,
			wantOK:       false,
			wantStatus:   http.StatusBadRequest,
			wantBodyIn:   "value is required",
		},
		{
			name:         "blank value when required",
			query:        "key=hello&value=%20%20",
			requireValue: true,
			wantOK:       false,
			wantStatus:   http.StatusBadRequest,
			wantBodyIn:   "value is required",
		},
		{
			name:         "value not required for get",
			query:        "key=hello",
			requireValue: false,
			wantOK:       true,
			wantKey:      "hello",
			wantValue:    "",
		},
		{
			name:         "key too long",
			query:        "key=" + strings.Repeat("a", maxKeyLen+1) + "&value=world",
			requireValue: true,
			wantOK:       false,
			wantStatus:   http.StatusBadRequest,
			wantBodyIn:   "key too long",
		},
		{
			name:         "key at max length is ok",
			query:        "key=" + strings.Repeat("a", maxKeyLen) + "&value=world",
			requireValue: true,
			wantOK:       true,
			wantKey:      strings.Repeat("a", maxKeyLen),
			wantValue:    "world",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/cache/?"+tc.query, nil)
			w := httptest.NewRecorder()

			key, value, ok := validateCacheInput(w, req, tc.requireValue)

			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if !ok {
				if w.Code != tc.wantStatus {
					t.Fatalf("status = %d, want %d", w.Code, tc.wantStatus)
				}
				if !strings.Contains(w.Body.String(), tc.wantBodyIn) {
					t.Fatalf("body %q does not contain %q", w.Body.String(), tc.wantBodyIn)
				}
				return
			}
			if key != tc.wantKey {
				t.Fatalf("key = %q, want %q", key, tc.wantKey)
			}
			if value != tc.wantValue {
				t.Fatalf("value = %q, want %q", value, tc.wantValue)
			}
		})
	}
}