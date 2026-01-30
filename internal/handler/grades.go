package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/q67q67q67-commits/college-system-app/internal/repository"
)

// GradesHandler — журнал оценок, GPA, транскрипт.
type GradesHandler struct{}

// ServeMyGrades: GET /api/grades — оценки текущего пользователя (студент).
func (h *GradesHandler) ServeMyGrades(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	userID := userIDFromRequest(r)
	list, err := repository.GradesByUserID(userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// ServeGPA: GET /api/grades/gpa — средний балл студента.
func (h *GradesHandler) ServeGPA(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	userID := userIDFromRequest(r)
	gpa, err := repository.GPAByUserID(userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]float64{"gpa": gpa})
}

// ServeTranscript: GET /api/grades/transcript — транскрипт (то же что оценки, сгруппировано по предметам при необходимости).
func (h *GradesHandler) ServeTranscript(w http.ResponseWriter, r *http.Request) {
	// Транскрипт = список оценок с датами и предметами (уже есть в GradesByUserID).
	h.ServeMyGrades(w, r)
}

// CreateGradeRequest — тело POST /api/grades (преподаватель).
type CreateGradeRequest struct {
	UserID     int64   `json:"user_id"`
	ScheduleID int64   `json:"schedule_id"`
	Grade      float64 `json:"grade"`
	GradeDate  string  `json:"grade_date"`
	Comment    string  `json:"comment"`
}

// ServeCreateGrade: POST /api/grades — добавить оценку (teacher).
func (h *GradesHandler) ServeCreateGrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var req CreateGradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if req.UserID == 0 || req.ScheduleID == 0 || req.GradeDate == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id, schedule_id, grade_date required"})
		return
	}
	id, err := repository.CreateGrade(req.UserID, req.ScheduleID, req.Grade, req.GradeDate, req.Comment)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

// ServeUpdateGrade: PUT /api/grades/:id — обновить оценку (teacher).
func (h *GradesHandler) ServeUpdateGrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "grade id required"})
		return
	}
	var req struct {
		Grade   float64 `json:"grade"`
		Comment string  `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if err := repository.UpdateGrade(id, req.Grade, req.Comment); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ServeDeleteGrade: DELETE /api/grades/:id — удалить оценку (teacher).
func (h *GradesHandler) ServeDeleteGrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "grade id required"})
		return
	}
	if err := repository.DeleteGrade(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
