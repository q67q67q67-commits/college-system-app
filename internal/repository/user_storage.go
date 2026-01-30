package repository

import (
	"database/sql"

	"github.com/q67q67q67-commits/college-system-app/internal/db"
)

const storageLimitBytes = 2147483648 // 2 ГБ

// UserStorageQuota возвращает лимит и использовано байт для пользователя.
func UserStorageQuota(userID int64) (limit, used int64, err error) {
	err = db.DB.QueryRow(`
		SELECT storage_limit_bytes, storage_used_bytes FROM users WHERE id = $1
	`, userID).Scan(&limit, &used)
	return limit, used, err
}

// CanAddStorage проверяет, можно ли добавить sizeBytes без превышения квоты.
func CanAddStorage(userID int64, sizeBytes int64) (bool, error) {
	limit, used, err := UserStorageQuota(userID)
	if err != nil {
		return false, err
	}
	return used+sizeBytes <= limit, nil
}

// AddUserFile добавляет запись о файле (триггер обновит storage_used_bytes).
func AddUserFile(userID int64, path, filename string, sizeBytes int64, contentType string) (int64, error) {
	var id int64
	err := db.DB.QueryRow(`
		INSERT INTO user_files (user_id, path, filename, size_bytes, content_type)
		VALUES ($1, $2, $3, $4, NULLIF($5,'')) RETURNING id
	`, userID, path, filename, sizeBytes, contentType).Scan(&id)
	return id, err
}

// DeleteUserFile удаляет запись о файле (триггер уменьшит storage_used_bytes).
func DeleteUserFile(id, userID int64) error {
	_, err := db.DB.Exec(`DELETE FROM user_files WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}

// UserFileRow — запись о файле пользователя.
type UserFileRow struct {
	ID          int64  `json:"id"`
	Path        string `json:"path"`
	Filename    string `json:"filename"`
	SizeBytes   int64  `json:"size_bytes"`
	ContentType string `json:"content_type"`
}

// ListUserFiles возвращает файлы пользователя.
func ListUserFiles(userID int64) ([]UserFileRow, error) {
	rows, err := db.DB.Query(`
		SELECT id, path, filename, size_bytes, COALESCE(content_type,'')
		FROM user_files WHERE user_id = $1 ORDER BY path
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []UserFileRow
	for rows.Next() {
		var f UserFileRow
		err := rows.Scan(&f.ID, &f.Path, &f.Filename, &f.SizeBytes, &f.ContentType)
		if err != nil {
			return nil, err
		}
		list = append(list, f)
	}
	return list, rows.Err()
}

// GetUserFilePath возвращает path файла по id и user_id (для скачивания/удаления).
func GetUserFilePath(id, userID int64) (path string, err error) {
	err = db.DB.QueryRow(`SELECT path FROM user_files WHERE id = $1 AND user_id = $2`, id, userID).Scan(&path)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return path, err
}

// GetUserFileForDownload возвращает path, filename, content_type для скачивания.
func GetUserFileForDownload(id, userID int64) (path, filename, contentType string, err error) {
	err = db.DB.QueryRow(`
		SELECT path, COALESCE(filename, path), COALESCE(content_type, 'application/octet-stream')
		FROM user_files WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(&path, &filename, &contentType)
	if err == sql.ErrNoRows {
		return "", "", "", nil
	}
	return path, filename, contentType, err
}
