package repository

import (
	"database/sql"
	"time"

	"github.com/narxoz-college/nc/internal/db"
)

// ForumPostRow — пост форума (для выдачи: автор скрыт при is_anonymous).
type ForumPostRow struct {
	ID          int64     `json:"id"`
	ParentID    *int64    `json:"parent_id"`
	AuthorID    *int64    `json:"author_id,omitempty"`   // не отдавать при is_anonymous
	AuthorName  string    `json:"author_name,omitempty"` // пусто при is_anonymous
	IsAnonymous bool      `json:"is_anonymous"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
}

// ForumPostsList возвращает посты: корневые (parent_id IS NULL) или ответы к parent_id.
func ForumPostsList(parentID *int64, limit, offset int) ([]ForumPostRow, error) {
	if limit <= 0 {
		limit = 20
	}
	var rows *sql.Rows
	var err error
	if parentID == nil {
		rows, err = db.DB.Query(`
			SELECT p.id, p.parent_id, p.author_id, u.full_name, p.is_anonymous,
			       COALESCE(p.title,''), p.body, p.created_at
			FROM forum_posts p
			LEFT JOIN users u ON u.id = p.author_id
			WHERE p.parent_id IS NULL
			ORDER BY p.created_at DESC LIMIT $1 OFFSET $2
		`, limit, offset)
	} else {
		rows, err = db.DB.Query(`
			SELECT p.id, p.parent_id, p.author_id, u.full_name, p.is_anonymous,
			       COALESCE(p.title,''), p.body, p.created_at
			FROM forum_posts p
			LEFT JOIN users u ON u.id = p.author_id
			WHERE p.parent_id = $1
			ORDER BY p.created_at
			LIMIT $2 OFFSET $3
		`, *parentID, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []ForumPostRow
	for rows.Next() {
		var post ForumPostRow
		var authorID sql.NullInt64
		var authorName sql.NullString
		var parentIDVal sql.NullInt64
		err := rows.Scan(&post.ID, &parentIDVal, &authorID, &authorName, &post.IsAnonymous,
			&post.Title, &post.Body, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		if parentIDVal.Valid {
			post.ParentID = &parentIDVal.Int64
		}
		if !post.IsAnonymous && authorID.Valid {
			post.AuthorID = &authorID.Int64
		}
		if !post.IsAnonymous && authorName.Valid {
			post.AuthorName = authorName.String
		}
		list = append(list, post)
	}
	return list, rows.Err()
}

// CreateForumPost создаёт пост (тему или ответ).
func CreateForumPost(parentID *int64, authorID int64, isAnonymous bool, title, body string) (int64, error) {
	var id int64
	err := db.DB.QueryRow(`
		INSERT INTO forum_posts (parent_id, author_id, is_anonymous, title, body)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`, parentID, authorID, isAnonymous, title, body).Scan(&id)
	return id, err
}
