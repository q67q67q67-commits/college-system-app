package handler

import (
	"net/http"
	"strconv"

	"github.com/q67q67q67-commits/college-system-app/internal/repository"
)

// GroupsHandler — список групп (преподаватель/админ).
type GroupsHandler struct{}

// ServeList: GET /api/groups — все группы с количеством студентов.
func (h *GroupsHandler) ServeList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	list, err := repository.ListGroups()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// ServeStudents: GET /api/groups/:id/students — студенты группы.
func (h *GroupsHandler) ServeStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	idStr := r.PathValue("id")
	groupID, _ := strconv.ParseInt(idStr, 10, 64)
	if groupID == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "group id required"})
		return
	}
	list, err := repository.StudentsByGroupID(groupID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, list)
}
