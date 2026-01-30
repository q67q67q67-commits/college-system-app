package handler

import (
	"encoding/json"
	"net/http"

	"github.com/q67q67q67-commits/college-system-app/internal/repository"
)

// BuildingMapHandler — карта здания.
type BuildingMapHandler struct{}

// ServeGet: GET /api/building-map — карта (все).
func (h *BuildingMapHandler) ServeGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	content, imageURL, err := repository.GetBuildingMap()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"content": content, "image_url": imageURL})
}

// ServeUpdate: PUT /api/building-map — обновить (admin/director).
func (h *BuildingMapHandler) ServeUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var req struct {
		Content  string `json:"content"`
		ImageURL string `json:"image_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if err := repository.UpsertBuildingMap(req.Content, req.ImageURL); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
