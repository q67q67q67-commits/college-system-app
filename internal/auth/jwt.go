package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const defaultExpire = 24 * time.Hour

// Claims — поля в JWT.
type Claims struct {
	jwt.RegisteredClaims
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
}

// CreateToken выдаёт JWT для пользователя.
func CreateToken(secret string, userID int64, role, email string, expire time.Duration) (string, error) {
	if expire == 0 {
		expire = defaultExpire
	}
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
		Role:   role,
		Email:  email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken проверяет JWT и возвращает userID, role, email.
func ValidateToken(secret, tokenString string) (userID int64, role, email string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, "", "", err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, "", "", errors.New("invalid token")
	}
	return claims.UserID, claims.Role, claims.Email, nil
}
