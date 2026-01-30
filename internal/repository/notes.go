package repository

import (
	"time"

	"github.com/q67q67q67-commits/college-system-app/internal/db"
)

// NoteRow — заметка пользователя.
type NoteRow struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// NoteCommentRow — комментарий к заметке.
type NoteCommentRow struct {
	ID        int64     `json:"id"`
	NoteID    int64     `json:"note_id"`
	UserID    int64     `json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// ListUserNotes возвращает заметки пользователя.
func ListUserNotes(userID int64, limit, offset int) ([]NoteRow, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := db.DB.Query(`
		SELECT id, user_id, COALESCE(title,''), body, created_at
		FROM user_notes WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []NoteRow
	for rows.Next() {
		var n NoteRow
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Body, &n.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, n)
	}
	return list, rows.Err()
}

// GetUserNote возвращает заметку по ID (только если принадлежит userID).
func GetUserNote(id, userID int64) (*NoteRow, error) {
	var n NoteRow
	err := db.DB.QueryRow(`
		SELECT id, user_id, COALESCE(title,''), body, created_at
		FROM user_notes WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(&n.ID, &n.UserID, &n.Title, &n.Body, &n.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

// CreateUserNote создаёт заметку.
func CreateUserNote(userID int64, title, body string) (int64, error) {
	var id int64
	err := db.DB.QueryRow(`
		INSERT INTO user_notes (user_id, title, body) VALUES ($1, $2, $3) RETURNING id
	`, userID, title, body).Scan(&id)
	return id, err
}

// ListNoteComments возвращает комментарии к заметке.
func ListNoteComments(noteID, ownerUserID int64) ([]NoteCommentRow, error) {
	rows, err := db.DB.Query(`
		SELECT c.id, c.note_id, c.user_id, c.body, c.created_at
		FROM user_note_comments c
		JOIN user_notes n ON n.id = c.note_id
		WHERE c.note_id = $1 AND n.user_id = $2
		ORDER BY c.created_at
	`, noteID, ownerUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []NoteCommentRow
	for rows.Next() {
		var c NoteCommentRow
		if err := rows.Scan(&c.ID, &c.NoteID, &c.UserID, &c.Body, &c.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

// CreateNoteComment добавляет комментарий к заметке (только владелец заметки).
func CreateNoteComment(noteID, userID int64, body string) (int64, error) {
	var id int64
	err := db.DB.QueryRow(`
		INSERT INTO user_note_comments (note_id, user_id, body)
		SELECT $1, $2, $3
		FROM user_notes WHERE id = $1 AND user_id = $2
		RETURNING id
	`, noteID, userID, body).Scan(&id)
	return id, err
}
