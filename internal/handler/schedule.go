package handler

import (
	"net/http"
	"strconv"

	"github.com/q67q67q67-commits/college-system-app/internal/repository"
)

// ScheduleHandler — расписание и дополнения к парам.
type ScheduleHandler struct{}

// ServeHTTP: GET /api/schedule — расписание текущего пользователя (студент: по группе; препод: по teacher_id).
func (h *ScheduleHandler) ServeSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	userID := userIDFromRequest(r)
	role := roleFromRequest(r)
	var list []repository.ScheduleRow
	var err error
	if role == "student" {
		groupID, err := repository.GroupIDByUserID(userID)
		if err != nil || groupID == 0 {
			writeJSON(w, http.StatusOK, []repository.ScheduleRow{})
			return
		}
		list, err = repository.SchedulesByGroupID(groupID)
	} else if role == "teacher" || role == "director" || role == "admin" {
		list, err = repository.SchedulesByTeacherID(userID)
	} else {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// ServeAttachments: GET /api/schedule/:id/attachments — дополнения к паре (ДЗ, уведомления).
func (h *ScheduleHandler) ServeAttachments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	idStr := r.PathValue("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "schedule id required"})
		return
	}
	scheduleID, _ := strconv.ParseInt(idStr, 10, 64)
	list, err := repository.AttachmentsByScheduleID(scheduleID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, list)
}
