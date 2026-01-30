package auth

import (
	"net/http"
	"strconv"
	"strings"
)

// RequireAuth проверяет наличие валидного JWT и кладёт userID, role, email в контекст запроса (заголовки для следующих handlers).
// Для WebSocket поддерживается токен в query: ?token=... (браузер не отправляет заголовки при апгрейде).
func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				if token := r.URL.Query().Get("token"); token != "" {
					auth = "Bearer " + token
					r.Header.Set("Authorization", auth)
				}
			}
			if auth == "" {
				http.Error(w, `{"error":"missing Authorization header"}`, http.StatusUnauthorized)
				return
			}
			const prefix = "Bearer "
			if !strings.HasPrefix(auth, prefix) {
				http.Error(w, `{"error":"invalid Authorization format"}`, http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimPrefix(auth, prefix)
			userID, role, email, fullName, err := ValidateToken(secret, tokenString)
			if err != nil {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}
			r.Header.Set("X-User-ID", strconv.FormatInt(userID, 10))
			r.Header.Set("X-User-Role", role)
			r.Header.Set("X-User-Email", email)
			r.Header.Set("X-User-Name", fullName)
			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole возвращает middleware, разрешающий только перечисленные роли.
func RequireRole(secret string, allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return RequireAuth(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Header.Get("X-User-Role")
			allowed := make(map[string]bool)
			for _, r := range allowedRoles {
				allowed[r] = true
			}
			if !allowed[role] {
				http.Error(w, `{"error":"forbidden: role not allowed"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}))
	}
}

