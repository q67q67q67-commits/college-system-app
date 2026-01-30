package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// userIDFromRequest возвращает user_id из заголовка (после middleware).
func userIDFromRequest(r *http.Request) int64 {
	s := r.Header.Get("X-User-ID")
	if s == "" {
		return 0
	}
	id, _ := strconv.ParseInt(s, 10, 64)
	return id
}

// roleFromRequest возвращает роль из заголовка.
func roleFromRequest(r *http.Request) string {
	return r.Header.Get("X-User-Role")
}

// writeJSON пишет JSON в ответ.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}
