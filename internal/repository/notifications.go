package repository

import (
	"database/sql"

	"github.com/q67q67q67-commits/college-system-app/internal/db"
)

// NotificationRow — уведомление.
type NotificationRow struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedBy int64  `json:"created_by"`
	CreatedAt string `json:"created_at"`
}

// ListNotifications возвращает последние уведомления.
func ListNotifications(limit int) ([]NotificationRow, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := db.DB.Query(`
		SELECT id, title, COALESCE(body,''), COALESCE(created_by,0), created_at::text
		FROM notifications ORDER BY created_at DESC LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []NotificationRow
	for rows.Next() {
		var n NotificationRow
		if err := rows.Scan(&n.ID, &n.Title, &n.Body, &n.CreatedBy, &n.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, n)
	}
	return list, rows.Err()
}

// CreateNotification создаёт уведомление.
func CreateNotification(title, body string, createdBy int64) (int64, error) {
	var id int64
	cb := sql.NullInt64{Int64: createdBy, Valid: createdBy > 0}
	err := db.DB.QueryRow(`
		INSERT INTO notifications (title, body, created_by) VALUES ($1, $2, $3) RETURNING id
	`, title, body, cb).Scan(&id)
	return id, err
}
