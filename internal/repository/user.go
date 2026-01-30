package repository

import (
	"database/sql"

	"github.com/q67q67q67-commits/college-system-app/internal/db"
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

// ProfileRow — данные профиля текущего пользователя.
type ProfileRow struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	AvatarURL string `json:"avatar_url"`
}

// GetProfileByID возвращает профиль пользователя по ID.
func GetProfileByID(userID int64) (*ProfileRow, error) {
	var p ProfileRow
	err := db.DB.QueryRow(`
		SELECT id, email, full_name, COALESCE(phone,''), role::text, COALESCE(avatar_url,'')
		FROM users WHERE id = $1 AND is_active = true
	`, userID).Scan(&p.ID, &p.Email, &p.FullName, &p.Phone, &p.Role, &p.AvatarURL)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
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

// UpdateUserAvatar обновляет avatar_url пользователя.
func UpdateUserAvatar(userID int64, avatarURL string) error {
	_, err := db.DB.Exec(`UPDATE users SET avatar_url = $1 WHERE id = $2`, avatarURL, userID)
	return err
}

// UserListItem — для списка всех пользователей (админ/директор).
type UserListItem struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Phone    string `json:"phone"`
}

// ListAllUsers возвращает всех пользователей (для админа/директора).
func ListAllUsers() ([]UserListItem, error) {
	rows, err := db.DB.Query(`
		SELECT id, full_name, email, role::text, COALESCE(phone,'')
		FROM users WHERE is_active = true ORDER BY full_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []UserListItem
	for rows.Next() {
		var u UserListItem
		if err := rows.Scan(&u.ID, &u.FullName, &u.Email, &u.Role, &u.Phone); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, rows.Err()
}

// UpdateUser обновляет данные пользователя (админ/директор).
func UpdateUser(id int64, fullName, email, phone, role string) error {
	_, err := db.DB.Exec(`
		UPDATE users SET
			full_name = CASE WHEN $1 <> '' THEN $1 ELSE full_name END,
			email = CASE WHEN $2 <> '' THEN $2 ELSE email END,
			phone = CASE WHEN $3 <> '' THEN $3 ELSE NULL END
		WHERE id = $4
	`, fullName, email, phone, id)
	if err != nil {
		return err
	}
	if role != "" && (role == "student" || role == "teacher" || role == "director" || role == "admin") {
		_, err = db.DB.Exec(`UPDATE users SET role = $1::user_role WHERE id = $2`, role, id)
	}
	return err
}

// DeactivateUser деактивирует пользователя (мягкое удаление).
func DeactivateUser(id int64) error {
	_, err := db.DB.Exec(`UPDATE users SET is_active = false WHERE id = $1`, id)
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
