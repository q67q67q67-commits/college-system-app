package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/q67q67q67-commits/college-system-app/internal/repository"
)

// LibraryHandler — поиск книг, экземпляры, бронирование.
type LibraryHandler struct{}

// ServeSearch: GET /api/library/books?q=... — поиск книг.
func (h *LibraryHandler) ServeSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	q := r.URL.Query().Get("q")
	if q == "" {
		q = " "
	}
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, _ := strconv.Atoi(l); n > 0 && n <= 100 {
			limit = n
		}
	}
	list, err := repository.SearchBooks(q, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// ServeCopies: GET /api/library/books/:id/copies — экземпляры книги (с владельцами).
func (h *LibraryHandler) ServeCopies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	idStr := r.PathValue("id")
	bookID, _ := strconv.ParseInt(idStr, 10, 64)
	if bookID == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "book id required"})
		return
	}
	list, err := repository.CopiesByBookID(bookID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// ReserveRequest — тело POST /api/library/reservations.
type ReserveRequest struct {
	BookID    int64  `json:"book_id"`
	ExpiresAt string `json:"expires_at"` // RFC3339 или дата
}

// ServeReserve: POST /api/library/reservations — забронировать книгу.
func (h *LibraryHandler) ServeReserve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	userID := userIDFromRequest(r)
	var req ReserveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if req.BookID == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "book_id required"})
		return
	}
	if req.ExpiresAt == "" {
		req.ExpiresAt = time.Now().Add(72 * time.Hour).Format(time.RFC3339)
	}
	id, err := repository.CreateReservation(req.BookID, userID, req.ExpiresAt)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}
