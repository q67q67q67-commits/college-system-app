package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/narxoz-college/nc/internal/repository"
)

// HomeworkHandler — управление ДЗ/уведомлениями к парам (teacher).
type HomeworkHandler struct{}

// CreateAttachmentRequest — тело POST /api/schedule/:id/attachments.
type CreateAttachmentRequest struct {
	Title          string `json:"title"`
	Body           string `json:"body"`
	AttachmentType string `json:"attachment_type"` // homework, notification, material
}

// ServeCreate: POST /api/schedule/:id/attachments — добавить ДЗ/уведомление (teacher).
func (h *HomeworkHandler) ServeCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	idStr := r.PathValue("id")
	scheduleID, _ := strconv.ParseInt(idStr, 10, 64)
	if scheduleID == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "schedule id required"})
		return
	}
	var req CreateAttachmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title required"})
		return
	}
	if req.AttachmentType == "" {
		req.AttachmentType = "homework"
	}
	id, err := repository.CreateAttachment(scheduleID, req.Title, req.Body, req.AttachmentType)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

// ServeUpdate: PUT /api/attachments/:id — обновить ДЗ (teacher).
func (h *HomeworkHandler) ServeUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
		return
	}
	var req struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if err := repository.UpdateAttachment(id, req.Title, req.Body); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ServeDelete: DELETE /api/attachments/:id — удалить ДЗ (teacher).
func (h *HomeworkHandler) ServeDelete(w http.ResponseWriter, r *http.Request) {
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
	if err := repository.DeleteAttachment(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
