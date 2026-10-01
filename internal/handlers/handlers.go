package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	PG  *pgxpool.Pool
	RDB *redis.Client
}

const maxKeyLen = 512

// validateCacheInput проверяет key и value для /cache/*.
// Возвращает (key, value, ok). Если ok=false - ответ уже отправлен.
func validateCacheInput(w http.ResponseWriter, r *http.Request, requireValue bool) (string, string, bool) {
	key := strings.TrimSpace(r.URL.Query().Get("key"))
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return "", "", false
	}
	if len(key) > maxKeyLen {
		http.Error(w, "key too long", http.StatusBadRequest)
		return "", "", false
	}

	value := r.URL.Query().Get("value")
	if requireValue && strings.TrimSpace(value) == "" {
		http.Error(w, "value is required", http.StatusBadRequest)
		return "", "", false
	}

	return key, value, true
}

func New(pg *pgxpool.Pool, rdb *redis.Client) *Handler {
	return &Handler{PG: pg, RDB: rdb}
}

// GET/POST /cache/set?key=foo&value=bar
func (h *Handler) CacheSet(w http.ResponseWriter, r *http.Request) {
	key, value, ok := validateCacheInput(w, r, true)
	if !ok {
		return
	}
	if err := h.RDB.Set(r.Context(), key, value, 0).Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"status": "ok", "key": key, "value": value})
}

// GET /cache/get?key=foo
func (h *Handler) CacheGet(w http.ResponseWriter, r *http.Request) {
	key, _, ok := validateCacheInput(w, r, false)
	if !ok {
		return
	}
	val, err := h.RDB.Get(r.Context(), key).Result()
	if err == redis.Nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"key": key, "value": val})
}

// GET /users
func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	rows, err := h.PG.Query(r.Context(), "SELECT id, name FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	users := []User{} 
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, users)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// GET /healthz 
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]string{"status": "ok"})
}

// GET /readyz 
func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := h.PG.Ping(ctx); err != nil {
		http.Error(w, "postgres not ready: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	if err := h.RDB.Ping(ctx).Err(); err != nil {
		http.Error(w, "redis not ready: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	writeJSON(w, map[string]string{"status": "ready"})
}