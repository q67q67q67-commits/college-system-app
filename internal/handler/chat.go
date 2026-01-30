package handler

import (
	"net/http"
	"strconv"

	"github.com/q67q67q67-commits/college-system-app/internal/repository"
)

// ChatHandler — REST API для чата (история, удаление).
type ChatHandler struct{}

// ServeMessages: GET /api/chat/messages — список сообщений.
func (h *ChatHandler) ServeMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, _ := strconv.Atoi(l); n > 0 && n <= 200 {
			limit = n
		}
	}
	list, err := repository.ListChatMessages(limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// ServeDeleteMessage: DELETE /api/chat/messages/:id — удалить сообщение.
func (h *ChatHandler) ServeDeleteMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
		return
	}
	userID := userIDFromRequest(r)
	role := roleFromRequest(r)
	isAdmin := role == "admin" || role == "director"
	ok, err := repository.DeleteChatMessage(id, userID, isAdmin)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if !ok {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
