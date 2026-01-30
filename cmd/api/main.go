// Narxoz College (NC) — точка входа API-сервера.
// Запуск: go run ./cmd/api
// Требует Go 1.22+ (PathValue в роутинге).
// Веб-прототип: http://localhost:8080/app/
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/q67q67q67-commits/college-system-app/internal/auth"
	"github.com/q67q67q67-commits/college-system-app/internal/chat"
	"github.com/q67q67q67-commits/college-system-app/internal/config"
	"github.com/q67q67q67-commits/college-system-app/internal/db"
	"github.com/q67q67q67-commits/college-system-app/internal/handler"
	"github.com/q67q67q67-commits/college-system-app/internal/repository"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatalf("api: %v", err)
	}
}

// corsMiddleware добавляет CORS-заголовки для веб-прототипа.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func run(args []string) error {
	cfg := config.Load()
	if err := db.Open(cfg.DBURL); err != nil {
		return err
	}
	defer db.Close()

	authHandler := &handler.AuthHandler{JWTSecret: cfg.JWTSecret}
	scheduleHandler := &handler.ScheduleHandler{}
	gradesHandler := &handler.GradesHandler{}
	groupsHandler := &handler.GroupsHandler{}
	libraryHandler := &handler.LibraryHandler{}
	forumHandler := &handler.ForumHandler{}
	eventsHandler := &handler.EventsHandler{}
	homeworkHandler := &handler.HomeworkHandler{}
	filesHandler := &handler.FilesHandler{UploadPath: cfg.UploadPath}
	profileHandler := &handler.ProfileHandler{UploadPath: cfg.UploadPath}
	directorHandler := &handler.DirectorHandler{}
	directorPostsHandler := &handler.DirectorPostsHandler{}
	notesHandler := &handler.NotesHandler{}
	buildingMapHandler := &handler.BuildingMapHandler{}
	usersHandler := &handler.UsersHandler{}
	notificationsHandler := &handler.NotificationsHandler{}
	chatHandler := &handler.ChatHandler{}

	chatHub := chat.NewHub()
	chatHub.SaveMessage = func(userID int64, fullName, body string) ([]byte, error) {
		id, author, err := repository.SaveChatMessage(userID, body, "")
		if err != nil {
			return nil, err
		}
		m := map[string]interface{}{"id": id, "user_id": userID, "author_name": author, "body": body}
		return json.Marshal(m)
	}
	go chatHub.Run()
	chatWS := chatHub.ServeWS(chat.GetUserIDFromHeader, chat.GetFullNameFromHeader)

	requireAuth := auth.RequireAuth(cfg.JWTSecret)
	requireTeacher := auth.RequireRole(cfg.JWTSecret, "teacher", "director", "admin")
	requireAdminDirector := auth.RequireRole(cfg.JWTSecret, "admin", "director")

	mux := http.NewServeMux()

	// Публичные
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("NC API. Prototype: /app/\n"))
	})
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	// Веб-прототип: статика из ./web по пути /app/
	webDir := filepath.Join("web")
	mux.Handle("GET /app/", http.StripPrefix("/app/", http.FileServer(http.Dir(webDir))))
	// Загруженные файлы (аватарки, посты директора)
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.UploadPath))))
	mux.HandleFunc("GET /app", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusFound)
	})

	// С авторизацией (любая роль)
	mux.Handle("GET /api/schedule", requireAuth(http.HandlerFunc(scheduleHandler.ServeSchedule)))
	mux.Handle("GET /api/schedule/{id}/attachments", requireAuth(http.HandlerFunc(scheduleHandler.ServeAttachments)))
	mux.Handle("GET /api/grades", requireAuth(http.HandlerFunc(gradesHandler.ServeMyGrades)))
	mux.Handle("GET /api/grades/gpa", requireAuth(http.HandlerFunc(gradesHandler.ServeGPA)))
	mux.Handle("GET /api/grades/transcript", requireAuth(http.HandlerFunc(gradesHandler.ServeTranscript)))
	mux.Handle("GET /api/groups", requireAuth(http.HandlerFunc(groupsHandler.ServeList)))
	mux.Handle("GET /api/groups/{id}/students", requireAuth(http.HandlerFunc(groupsHandler.ServeStudents)))
	mux.Handle("GET /api/library/books", requireAuth(http.HandlerFunc(libraryHandler.ServeSearch)))
	mux.Handle("GET /api/library/books/{id}/copies", requireAuth(http.HandlerFunc(libraryHandler.ServeCopies)))
	mux.Handle("POST /api/library/reservations", requireAuth(http.HandlerFunc(libraryHandler.ServeReserve)))
	mux.Handle("GET /api/forum/posts", requireAuth(http.HandlerFunc(forumHandler.ServeList)))
	mux.Handle("POST /api/forum/posts", requireAuth(http.HandlerFunc(forumHandler.ServeCreate)))
	mux.Handle("DELETE /api/forum/posts/{id}", requireAuth(http.HandlerFunc(forumHandler.ServeDelete)))
	mux.Handle("GET /api/events", requireAuth(http.HandlerFunc(eventsHandler.ServeList)))
	mux.Handle("GET /api/events/{id}", requireAuth(http.HandlerFunc(eventsHandler.ServeGet)))
	mux.Handle("GET /api/files", requireAuth(http.HandlerFunc(filesHandler.ServeList)))
	mux.Handle("GET /api/files/{id}/download", requireAuth(http.HandlerFunc(filesHandler.ServeDownload)))
	mux.Handle("POST /api/files/upload", requireAuth(http.HandlerFunc(filesHandler.ServeUpload)))
	mux.Handle("DELETE /api/files/{id}", requireAuth(http.HandlerFunc(filesHandler.ServeDelete)))
	mux.Handle("GET /api/chat/ws", requireAuth(http.HandlerFunc(chatWS)))
	mux.Handle("GET /api/chat/messages", requireAuth(http.HandlerFunc(chatHandler.ServeMessages)))
	mux.Handle("DELETE /api/chat/messages/{id}", requireAuth(http.HandlerFunc(chatHandler.ServeDeleteMessage)))
	mux.Handle("GET /api/profile", requireAuth(http.HandlerFunc(profileHandler.ServeGet)))
	mux.Handle("PUT /api/profile", requireAuth(http.HandlerFunc(profileHandler.ServeUpdate)))
	mux.Handle("POST /api/profile/avatar", requireAuth(http.HandlerFunc(profileHandler.ServeAvatar)))
	mux.Handle("GET /api/director", http.HandlerFunc(directorHandler.ServeGet))
	mux.Handle("GET /api/director/posts", http.HandlerFunc(directorPostsHandler.ServeList))
	mux.Handle("GET /api/director/posts/{id}", requireAuth(http.HandlerFunc(directorPostsHandler.ServeGet)))
	mux.Handle("POST /api/director/posts/{id}/comments", requireAuth(http.HandlerFunc(directorPostsHandler.ServeAddComment)))
	mux.Handle("GET /api/building-map", http.HandlerFunc(buildingMapHandler.ServeGet))
	mux.Handle("GET /api/notes", requireAuth(http.HandlerFunc(notesHandler.ServeList)))
	mux.Handle("POST /api/notes", requireAuth(http.HandlerFunc(notesHandler.ServeCreate)))
	mux.Handle("GET /api/notes/{id}", requireAuth(http.HandlerFunc(notesHandler.ServeGet)))
	mux.Handle("POST /api/notes/{id}/comments", requireAuth(http.HandlerFunc(notesHandler.ServeAddComment)))

	// Только преподаватель: оценки, ДЗ
	mux.Handle("POST /api/grades", requireTeacher(http.HandlerFunc(gradesHandler.ServeCreateGrade)))
	mux.Handle("PUT /api/grades/{id}", requireTeacher(http.HandlerFunc(gradesHandler.ServeUpdateGrade)))
	mux.Handle("DELETE /api/grades/{id}", requireTeacher(http.HandlerFunc(gradesHandler.ServeDeleteGrade)))
	mux.Handle("POST /api/schedule/{id}/attachments", requireTeacher(http.HandlerFunc(homeworkHandler.ServeCreate)))
	mux.Handle("PUT /api/attachments/{id}", requireTeacher(http.HandlerFunc(homeworkHandler.ServeUpdate)))
	mux.Handle("DELETE /api/attachments/{id}", requireTeacher(http.HandlerFunc(homeworkHandler.ServeDelete)))

	// Только админ/директор: CRUD событий, пользователи, уведомления
	mux.Handle("POST /api/events", requireAdminDirector(http.HandlerFunc(eventsHandler.ServeCreate)))
	mux.Handle("PUT /api/events/{id}", requireAdminDirector(http.HandlerFunc(eventsHandler.ServeUpdate)))
	mux.Handle("DELETE /api/events/{id}", requireAdminDirector(http.HandlerFunc(eventsHandler.ServeDelete)))
	mux.Handle("GET /api/users", requireAdminDirector(http.HandlerFunc(usersHandler.ServeList)))
	mux.Handle("GET /api/users/{id}", requireAuth(http.HandlerFunc(usersHandler.ServeGet)))
	mux.Handle("PUT /api/users/{id}", requireAdminDirector(http.HandlerFunc(usersHandler.ServeUpdate)))
	mux.Handle("DELETE /api/users/{id}", requireAdminDirector(http.HandlerFunc(usersHandler.ServeDelete)))
	mux.Handle("GET /api/notifications", requireAuth(http.HandlerFunc(notificationsHandler.ServeList)))
	mux.Handle("POST /api/notifications", requireAdminDirector(http.HandlerFunc(notificationsHandler.ServeCreate)))
	mux.Handle("POST /api/director/posts", requireAuth(http.HandlerFunc(directorPostsHandler.ServeCreate)))
	mux.Handle("PUT /api/building-map", requireAdminDirector(http.HandlerFunc(buildingMapHandler.ServeUpdate)))

	log.Printf("api: listening on %s", cfg.Addr)
	return http.ListenAndServe(cfg.Addr, corsMiddleware(mux))
}
