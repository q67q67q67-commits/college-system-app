package repository

import (
	"time"

	"github.com/q67q67q67-commits/college-system-app/internal/db"
)

// DirectorPostRow — пост директора (как в форуме: title, body).
type DirectorPostRow struct {
	ID        int64     `json:"id"`
	DirectorID int64    `json:"director_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// DirectorPostCommentRow — комментарий к посту директора.
type DirectorPostCommentRow struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	AuthorName string   `json:"author_name"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// ListDirectorPosts возвращает посты директора.
func ListDirectorPosts(directorID int64) ([]DirectorPostRow, error) {
	rows, err := db.DB.Query(`
		SELECT id, director_id, COALESCE(title,''), COALESCE(body,''), created_at
		FROM director_posts WHERE director_id = $1 ORDER BY created_at DESC
	`, directorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []DirectorPostRow
	for rows.Next() {
		var p DirectorPostRow
		if err := rows.Scan(&p.ID, &p.DirectorID, &p.Title, &p.Body, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

// GetDirectorPost возвращает пост по ID.
func GetDirectorPost(id int64) (*DirectorPostRow, error) {
	var p DirectorPostRow
	err := db.DB.QueryRow(`
		SELECT id, director_id, COALESCE(title,''), COALESCE(body,''), created_at
		FROM director_posts WHERE id = $1
	`, id).Scan(&p.ID, &p.DirectorID, &p.Title, &p.Body, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// CreateDirectorPost создаёт пост (title, body).
func CreateDirectorPost(directorID int64, title, body string) (int64, error) {
	var id int64
	err := db.DB.QueryRow(`
		INSERT INTO director_posts (director_id, title, body) VALUES ($1, $2, $3) RETURNING id
	`, directorID, title, body).Scan(&id)
	return id, err
}

// ListDirectorPostComments возвращает комментарии к посту.
func ListDirectorPostComments(postID int64) ([]DirectorPostCommentRow, error) {
	rows, err := db.DB.Query(`
		SELECT c.id, c.post_id, c.user_id, COALESCE(u.full_name,''), c.body, c.created_at
		FROM director_post_comments c
		LEFT JOIN users u ON u.id = c.user_id
		WHERE c.post_id = $1 ORDER BY c.created_at
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []DirectorPostCommentRow
	for rows.Next() {
		var c DirectorPostCommentRow
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.AuthorName, &c.Body, &c.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

// CreateDirectorPostComment добавляет комментарий.
func CreateDirectorPostComment(postID, userID int64, body string) (int64, error) {
	var id int64
	err := db.DB.QueryRow(`
		INSERT INTO director_post_comments (post_id, user_id, body) VALUES ($1, $2, $3) RETURNING id
	`, postID, userID, body).Scan(&id)
	return id, err
}
