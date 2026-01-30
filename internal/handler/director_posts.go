package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/q67q67q67-commits/college-system-app/internal/repository"
)

// DirectorPostsHandler — посты директора (как форум: title, body, комментарии).
type DirectorPostsHandler struct{}

// ServeList: GET /api/director/posts — посты директора.
func (h *DirectorPostsHandler) ServeList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	d, err := repository.DirectorProfile()
	if err != nil || d == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	list, err := repository.ListDirectorPosts(d.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// ServeCreate: POST /api/director/posts — добавить пост (только директор).
func (h *DirectorPostsHandler) ServeCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	userID := userIDFromRequest(r)
	role := roleFromRequest(r)
	if role != "director" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "only director can add posts"})
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
	if req.Body == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body required"})
		return
	}
	id, err := repository.CreateDirectorPost(userID, req.Title, req.Body)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

// ServeGet: GET /api/director/posts/:id — пост с комментариями.
func (h *DirectorPostsHandler) ServeGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
		return
	}
	post, err := repository.GetDirectorPost(id)
	if err != nil || post == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	comments, _ := repository.ListDirectorPostComments(id)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"post":     post,
		"comments": comments,
	})
}

// ServeAddComment: POST /api/director/posts/:id/comments — добавить комментарий.
func (h *DirectorPostsHandler) ServeAddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	userID := userIDFromRequest(r)
	if userID == 0 {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	idStr := r.PathValue("id")
	postID, _ := strconv.ParseInt(idStr, 10, 64)
	if postID == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
		return
	}
	var req struct {
		Body string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if req.Body == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body required"})
		return
	}
	id, err := repository.CreateDirectorPostComment(postID, userID, req.Body)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}
