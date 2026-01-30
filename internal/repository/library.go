package repository

import (
	"database/sql"
	"time"

	"github.com/narxoz-college/nc/internal/db"
)

// LibraryBookRow — книга.
type LibraryBookRow struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	ISBN        string `json:"isbn"`
	Description string `json:"description"`
	HasPhysical bool   `json:"has_physical"`
	TotalCopies int    `json:"total_copies"`
	Available   int    `json:"available"`
}

// LibraryCopyRow — экземпляр с владельцем.
type LibraryCopyRow struct {
	ID         int64     `json:"id"`
	BookID     int64     `json:"book_id"`
	HolderID   *int64    `json:"holder_id"`
	HolderName string    `json:"holder_name,omitempty"`
	BorrowedAt *time.Time `json:"borrowed_at,omitempty"`
	DueAt      *time.Time `json:"due_at,omitempty"`
}

// SearchBooks ищет книги по title/author (подстрока).
func SearchBooks(query string, limit int) ([]LibraryBookRow, error) {
	if limit <= 0 {
		limit = 20
	}
	q := "%" + query + "%"
	rows, err := db.DB.Query(`
		SELECT b.id, b.title, COALESCE(b.author,''), COALESCE(b.isbn,''), COALESCE(b.description,''),
		       b.has_physical, b.total_copies,
		       b.total_copies - (SELECT COUNT(*) FROM library_copies c WHERE c.book_id = b.id AND c.holder_id IS NOT NULL)
		FROM library_books b
		WHERE b.title ILIKE $1 OR b.author ILIKE $1
		ORDER BY b.title LIMIT $2
	`, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []LibraryBookRow
	for rows.Next() {
		var book LibraryBookRow
		var avail int
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Description,
			&book.HasPhysical, &book.TotalCopies, &avail)
		if err != nil {
			return nil, err
		}
		book.Available = avail
		list = append(list, book)
	}
	return list, rows.Err()
}

// CopiesByBookID возвращает экземпляры книги с владельцами.
func CopiesByBookID(bookID int64) ([]LibraryCopyRow, error) {
	rows, err := db.DB.Query(`
		SELECT c.id, c.book_id, c.holder_id, u.full_name, c.borrowed_at, c.due_at
		FROM library_copies c
		LEFT JOIN users u ON u.id = c.holder_id
		WHERE c.book_id = $1 ORDER BY c.id
	`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []LibraryCopyRow
	for rows.Next() {
		var copy LibraryCopyRow
		var holderID sql.NullInt64
		var holderName sql.NullString
		err := rows.Scan(&copy.ID, &copy.BookID, &holderID, &holderName, &copy.BorrowedAt, &copy.DueAt)
		if err != nil {
			return nil, err
		}
		if holderID.Valid {
			copy.HolderID = &holderID.Int64
		}
		if holderName.Valid {
			copy.HolderName = holderName.String
		}
		list = append(list, copy)
	}
	return list, rows.Err()
}

// CreateReservation создаёт бронь книги.
func CreateReservation(bookID, userID int64, expiresAt string) (int64, error) {
	var id int64
	err := db.DB.QueryRow(`
		INSERT INTO library_reservations (book_id, user_id, expires_at)
		VALUES ($1, $2, $3::timestamptz) RETURNING id
	`, bookID, userID, expiresAt).Scan(&id)
	return id, err
}
