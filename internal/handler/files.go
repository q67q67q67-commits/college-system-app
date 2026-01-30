package handler

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/narxoz-college/nc/internal/repository"
)

// FilesHandler — загрузка/список/удаление файлов с проверкой квоты 2 ГБ.
type FilesHandler struct {
	UploadPath string
}

// ServeList: GET /api/files — список файлов пользователя.
func (h *FilesHandler) ServeList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	userID := userIDFromRequest(r)
	list, err := repository.ListUserFiles(userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	limit, used, _ := repository.UserStorageQuota(userID)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"files":           list,
		"storage_limit":   limit,
		"storage_used":    used,
	})
}

// ServeUpload: POST /api/files/upload — загрузить файл (multipart/form-data, поле "file"). Проверка квоты 2 ГБ.
func (h *FilesHandler) ServeUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	userID := userIDFromRequest(r)
	_ = r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "file required"})
		return
	}
	defer file.Close()
	size := header.Size
	ok, err := repository.CanAddStorage(userID, size)
	if err != nil || !ok {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "storage quota exceeded (2 GB limit)"})
		return
	}
	filename := header.Filename
	if filename == "" {
		filename = "file"
	}
	relPath := filepath.Join(strconv.FormatInt(userID, 10), filename)
	fullPath := filepath.Join(h.UploadPath, relPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	dst, err := os.Create(fullPath)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	defer dst.Close()
	n, err := io.Copy(dst, file)
	if err != nil {
		os.Remove(fullPath)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	contentType := header.Header.Get("Content-Type")
	id, err := repository.AddUserFile(userID, relPath, filename, n, contentType)
	if err != nil {
		os.Remove(fullPath)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id":   id,
		"path": relPath,
		"size": n,
	})
}

// ServeDelete: DELETE /api/files/:id — удалить файл пользователя.
func (h *FilesHandler) ServeDelete(w http.ResponseWriter, r *http.Request) {
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
	path, err := repository.GetUserFilePath(id, userID)
	if err != nil || path == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	fullPath := filepath.Join(h.UploadPath, path)
	_ = os.Remove(fullPath)
	if err := repository.DeleteUserFile(id, userID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

