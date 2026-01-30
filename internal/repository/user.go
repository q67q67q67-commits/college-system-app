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

// DirectorProfileRow — публичный профиль директора (для страницы «Профиль директора»).
type DirectorProfileRow struct {
	ID        int64  `json:"id"`
	FullName  string `json:"full_name"`
	AvatarURL string `json:"avatar_url"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
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

// UpdateUserPassword обновляет пароль пользователя (хэш должен быть передан уже захэшированным).
func UpdateUserPassword(userID int64, passwordHash string) error {
	_, err := db.DB.Exec(`UPDATE users SET password_hash = $1 WHERE id = $2`, passwordHash, userID)
	return err
}

// UpdateUserPhone обновляет телефон пользователя.
func UpdateUserPhone(userID int64, phone string) error {
	_, err := db.DB.Exec(`UPDATE users SET phone = $1 WHERE id = $2`, phone, userID)
	return err
}

// DirectorProfile возвращает публичный профиль директора (первый пользователь с ролью director).
func DirectorProfile() (*DirectorProfileRow, error) {
	var d DirectorProfileRow
	err := db.DB.QueryRow(`
		SELECT id, full_name, COALESCE(avatar_url,''), COALESCE(phone,''), COALESCE(email,'')
		FROM users WHERE role = 'director' AND is_active = true ORDER BY id LIMIT 1
	`).Scan(&d.ID, &d.FullName, &d.AvatarURL, &d.Phone, &d.Email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}
