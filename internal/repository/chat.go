package repository

import (
	"database/sql"

	"github.com/q67q67q67-commits/college-system-app/internal/db"
)

// ChatMessageRow — сообщение чата.
type ChatMessageRow struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Author    string `json:"author_name"`
	Body      string `json:"body"`
	MediaURL  string `json:"media_url,omitempty"`
	CreatedAt string `json:"created_at"`
}

// SaveChatMessage сохраняет сообщение и возвращает id.
func SaveChatMessage(userID int64, body, mediaURL string) (int64, string, error) {
	var id int64
	var author string
	err := db.DB.QueryRow(`
		INSERT INTO chat_messages (user_id, body, media_url)
		SELECT $1, $2, NULLIF($3,'')
		RETURNING id, (SELECT full_name FROM users WHERE id = $1)
	`, userID, body, mediaURL).Scan(&id, &author)
	return id, author, err
}

// ListChatMessages возвращает последние сообщения.
func ListChatMessages(limit int) ([]ChatMessageRow, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := db.DB.Query(`
		SELECT m.id, m.user_id, u.full_name, m.body, COALESCE(m.media_url,''), m.created_at::text
		FROM chat_messages m
		JOIN users u ON u.id = m.user_id
		ORDER BY m.created_at DESC LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []ChatMessageRow
	for rows.Next() {
		var m ChatMessageRow
		if err := rows.Scan(&m.ID, &m.UserID, &m.Author, &m.Body, &m.MediaURL, &m.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	// reverse to chronological
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list, rows.Err()
}

// DeleteChatMessage удаляет сообщение. Возвращает true если удалено.
func DeleteChatMessage(id, userID int64, isAdmin bool) (bool, error) {
	var res int64
	if isAdmin {
		err := db.DB.QueryRow(`DELETE FROM chat_messages WHERE id = $1 RETURNING id`, id).Scan(&res)
		return err == nil, err
	}
	err := db.DB.QueryRow(`DELETE FROM chat_messages WHERE id = $1 AND user_id = $2 RETURNING id`, id, userID).Scan(&res)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}
