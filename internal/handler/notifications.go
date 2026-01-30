package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/q67q67q67-commits/college-system-app/internal/repository"
)

// NotificationsHandler — уведомления.
type NotificationsHandler struct{}

// ServeList: GET /api/notifications — список уведомлений.
func (h *NotificationsHandler) ServeList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, _ := strconv.Atoi(l); n > 0 && n <= 100 {
			limit = n
		}
	}
	list, err := repository.ListNotifications(limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// CreateNotificationRequest — тело POST /api/notifications.
type CreateNotificationRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// ServeCreate: POST /api/notifications — создать уведомление (админ/директор).
func (h *NotificationsHandler) ServeCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	userID := userIDFromRequest(r)
	var req CreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title required"})
		return
	}
	id, err := repository.CreateNotification(req.Title, req.Body, userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}
