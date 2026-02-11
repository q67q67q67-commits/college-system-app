package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/narxoz-college/nc/internal/auth"
	"github.com/narxoz-college/nc/internal/repository"
)

// LoginRequest — тело POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse — ответ при успешном логине.
type LoginResponse struct {
	Token   string `json:"token"`
	UserID  int64  `json:"user_id"`
	Role    string `json:"role"`
	Email   string `json:"email"`
	FullName string `json:"full_name"`
	Expires string `json:"expires_at"`
}

// AuthHandler обрабатывает POST /auth/login.
type AuthHandler struct {
	JWTSecret string
}

// Login — POST /auth/login: проверяет email+password, возвращает JWT.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error":"email and password required"}`, http.StatusBadRequest)
		return
	}
	u, err := repository.UserByEmail(req.Email)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	if u == nil {
		http.Error(w, `{"error":"invalid email or password"}`, http.StatusUnauthorized)
		return
	}
	// Проверка доступа после ухода: если left_at задан и прошло > 1 года — вход запрещён.
	if u.LeftAt.Valid {
		if time.Since(u.LeftAt.Time) > 365*24*time.Hour {
			http.Error(w, `{"error":"access expired after leaving college"}`, http.StatusForbidden)
			return
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, `{"error":"invalid email or password"}`, http.StatusUnauthorized)
		return
	}
	expire := 24 * time.Hour
	token, err := auth.CreateToken(h.JWTSecret, u.ID, u.Role, u.Email, expire)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Token:    token,
		UserID:   u.ID,
		Role:     u.Role,
		Email:    u.Email,
		FullName: u.FullName,
		Expires:  time.Now().Add(expire).Format(time.RFC3339),
	})
}
