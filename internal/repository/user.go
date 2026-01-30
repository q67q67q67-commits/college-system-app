package repository

import (
	"database/sql"

	"github.com/narxoz-college/nc/internal/db"
)

// UserRow — строка users для логина.
type UserRow struct {
	ID           int64
	Email        string
	PasswordHash string
	Role         string
	FullName     string
	IsActive     bool
}

// UserByEmail возвращает пользователя по email.
func UserByEmail(email string) (*UserRow, error) {
	var u UserRow
	err := db.DB.QueryRow(`
		SELECT id, email, password_hash, role::text, full_name, is_active
		FROM users WHERE email = $1 AND is_active = true
	`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.FullName, &u.IsActive)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
