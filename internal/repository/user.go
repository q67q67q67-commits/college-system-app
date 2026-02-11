package repository

import (
	"database/sql"
	"time"

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
	LeftAt       sql.NullTime
}

// UserByEmail возвращает пользователя по email.
func UserByEmail(email string) (*UserRow, error) {
	var u UserRow
	err := db.DB.QueryRow(`
		SELECT id, email, password_hash, role::text, full_name, is_active, left_at
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

// UserProfile — расширенные данные профиля для текущего пользователя.
type UserProfile struct {
	ID                int64     `json:"id"`
	Email             string    `json:"email"`
	Role              string    `json:"role"`
	FullName          string    `json:"full_name"`
	Phone             string    `json:"phone"`
	AvatarURL         string    `json:"avatar_url"`
	Language          string    `json:"language"`
	StorageLimitBytes int64     `json:"storage_limit_bytes"`
	StorageUsedBytes  int64     `json:"storage_used_bytes"`
	LeftAt            *time.Time `json:"left_at,omitempty"`
}

// UserProfileByID возвращает профиль пользователя по его id.
func UserProfileByID(userID int64) (*UserProfile, error) {
	var p UserProfile
	var leftAt sql.NullTime
	err := db.DB.QueryRow(`
		SELECT id, email, role::text, full_name,
		       COALESCE(phone, ''), COALESCE(avatar_url, ''), language,
		       storage_limit_bytes, storage_used_bytes, left_at
		FROM users
		WHERE id = $1
	`, userID).Scan(
		&p.ID, &p.Email, &p.Role, &p.FullName,
		&p.Phone, &p.AvatarURL, &p.Language,
		&p.StorageLimitBytes, &p.StorageUsedBytes, &leftAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if leftAt.Valid {
		t := leftAt.Time
		p.LeftAt = &t
	}
	return &p, nil
}

// UpdateUserProfile обновляет базовые поля профиля пользователя.
func UpdateUserProfile(userID int64, fullName, phone, language, avatarURL string) error {
	_, err := db.DB.Exec(`
		UPDATE users
		SET full_name = $1,
		    phone      = NULLIF($2, ''),
		    language   = $3,
		    avatar_url = NULLIF($4, '')
		WHERE id = $5
	`, fullName, phone, language, avatarURL, userID)
	return err
}

// UpdateUserPassword обновляет хеш пароля пользователя.
func UpdateUserPassword(userID int64, newHash string) error {
	_, err := db.DB.Exec(`UPDATE users SET password_hash = $1 WHERE id = $2`, newHash, userID)
	return err
}

// UserPasswordByID возвращает email и хеш пароля пользователя по id (для смены пароля).
func UserPasswordByID(userID int64) (*UserRow, error) {
	var u UserRow
	err := db.DB.QueryRow(`
		SELECT id, email, password_hash, role::text, full_name, is_active, left_at
		FROM users WHERE id = $1
	`, userID).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.FullName, &u.IsActive, &u.LeftAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
