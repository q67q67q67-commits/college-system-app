package handler

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/narxoz-college/nc/internal/repository"
)

// ProfileHandler — обновление профиля (пароль, телефон).
type ProfileHandler struct{}

// UpdateProfileRequest — тело PUT /api/profile. Передавать только те поля, которые нужно обновить.
type UpdateProfileRequest struct {
	Password *string `json:"password"`
	Phone    *string `json:"phone"`
}

// ServeUpdate: PUT /api/profile — обновить пароль и/или телефон текущего пользователя.
func (h *ProfileHandler) ServeUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	userID := userIDFromRequest(r)
	if userID == 0 {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if req.Password != nil && *req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
			return
		}
		if err := repository.UpdateUserPassword(userID, string(hash)); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
			return
		}
	}
	if req.Phone != nil {
		phone := *req.Phone
		if err := repository.UpdateUserPhone(userID, phone); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// DirectorHandler — публичный профиль директора.
type DirectorHandler struct{}

// ServeGet: GET /api/director — публичный профиль директора (для страницы «Профиль директора»).
func (h *DirectorHandler) ServeGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	d, err := repository.DirectorProfile()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if d == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, d)
}
