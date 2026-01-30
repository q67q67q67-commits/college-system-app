package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/q67q67q67-commits/college-system-app/internal/repository"
)

// ForumHandler — форум (темы, ответы, анонимные посты).
type ForumHandler struct{}

// ServeList: GET /api/forum/posts?parent_id=...&limit=...&offset=... — список постов (корневые или ответы).
func (h *ForumHandler) ServeList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var parentID *int64
	if p := r.URL.Query().Get("parent_id"); p != "" {
		id, _ := strconv.ParseInt(p, 10, 64)
		parentID = &id
	}
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, _ := strconv.Atoi(l); n > 0 && n <= 100 {
			limit = n
		}
	}
	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if n, _ := strconv.Atoi(o); n >= 0 {
			offset = n
		}
	}
	list, err := repository.ForumPostsList(parentID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// CreateForumPostRequest — тело POST /api/forum/posts.
type CreateForumPostRequest struct {
	ParentID    *int64 `json:"parent_id"`
	IsAnonymous bool   `json:"is_anonymous"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	MediaURL    string `json:"media_url"`
}

// ServeCreate: POST /api/forum/posts — создать тему или ответ (с опцией анонимно).
func (h *ForumHandler) ServeCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	userID := userIDFromRequest(r)
	var req CreateForumPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if req.Body == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body required"})
		return
	}
	id, err := repository.CreateForumPost(req.ParentID, userID, req.IsAnonymous, req.Title, req.Body, req.MediaURL)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

// ServeDelete: DELETE /api/forum/posts/:id — удалить пост (свой или любой для admin/director).
func (h *ForumHandler) ServeDelete(w http.ResponseWriter, r *http.Request) {
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
	ok, err := repository.DeleteForumPost(id, userID, isAdmin)
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
