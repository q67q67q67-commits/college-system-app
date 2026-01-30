package repository

import (
	"database/sql"
	"time"

	"github.com/narxoz-college/nc/internal/db"
)

// EventRow — событие/новость.
type EventRow struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	EventDate   time.Time `json:"event_date"`
	Location    string    `json:"location"`
	CreatedBy   *int64    `json:"created_by,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// EventsList возвращает события с опциональной фильтрацией по дате и поиску.
func EventsList(fromDate, toDate, search string, limit, offset int) ([]EventRow, error) {
	if limit <= 0 {
		limit = 20
	}
	query := `
		SELECT id, title, COALESCE(description,''), COALESCE(image_url,''), event_date,
		       COALESCE(location,''), created_by, created_at
		FROM events
		WHERE event_date >= COALESCE(NULLIF($1,'')::date, '0001-01-01'::date)
		  AND event_date <= COALESCE(NULLIF($2,'')::date, '9999-12-31'::date)
		  AND ($3 = '' OR title ILIKE $3 OR description ILIKE $3)
		ORDER BY event_date DESC LIMIT $4 OFFSET $5
	`
	searchLike := ""
	if search != "" {
		searchLike = "%" + search + "%"
	}
	rows, err := db.DB.Query(query, fromDate, toDate, searchLike, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []EventRow
	for rows.Next() {
		var e EventRow
		var createdBy sql.NullInt64
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.ImageURL, &e.EventDate, &e.Location, &createdBy, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		if createdBy.Valid {
			e.CreatedBy = &createdBy.Int64
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

// EventByID возвращает событие по ID.
func EventByID(id int64) (*EventRow, error) {
	var e EventRow
	var createdBy sql.NullInt64
	err := db.DB.QueryRow(`
		SELECT id, title, COALESCE(description,''), COALESCE(image_url,''), event_date,
		       COALESCE(location,''), created_by, created_at
		FROM events WHERE id = $1
	`, id).Scan(&e.ID, &e.Title, &e.Description, &e.ImageURL, &e.EventDate, &e.Location, &createdBy, &e.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if createdBy.Valid {
		e.CreatedBy = &createdBy.Int64
	}
	return &e, nil
}

// CreateEvent создаёт событие.
func CreateEvent(title, description, imageURL, eventDate, location string, createdBy int64) (int64, error) {
	var id int64
	err := db.DB.QueryRow(`
		INSERT INTO events (title, description, image_url, event_date, location, created_by)
		VALUES ($1, $2, NULLIF($3,''), $4::timestamptz, NULLIF($5,''), $6) RETURNING id
	`, title, description, imageURL, eventDate, location, createdBy).Scan(&id)
	return id, err
}

// UpdateEvent обновляет событие.
func UpdateEvent(id int64, title, description, imageURL, eventDate, location string) error {
	_, err := db.DB.Exec(`
		UPDATE events SET title = $1, description = $2, image_url = NULLIF($3,''),
		                  event_date = $4::timestamptz, location = NULLIF($5,'')
		WHERE id = $6
	`, title, description, imageURL, eventDate, location, id)
	return err
}

// DeleteEvent удаляет событие.
func DeleteEvent(id int64) error {
	_, err := db.DB.Exec(`DELETE FROM events WHERE id = $1`, id)
	return err
}
