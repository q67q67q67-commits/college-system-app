// Narxoz College (NC) — точка входа API-сервера.
// Запуск: go run ./cmd/api
// Требует Go 1.22+ (PathValue в роутинге).
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/narxoz-college/nc/internal/auth"
	"github.com/narxoz-college/nc/internal/chat"
	"github.com/narxoz-college/nc/internal/config"
	"github.com/narxoz-college/nc/internal/db"
	"github.com/narxoz-college/nc/internal/handler"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatalf("api: %v", err)
	}
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

	chatHub := chat.NewHub()
	go chatHub.Run()
	chatWS := chatHub.ServeWS(chat.GetUserIDFromHeader)

	requireAuth := auth.RequireAuth(cfg.JWTSecret)
	requireTeacher := auth.RequireRole(cfg.JWTSecret, "teacher", "director", "admin")
	requireAdminDirector := auth.RequireRole(cfg.JWTSecret, "admin", "director")

	mux := http.NewServeMux()

	// Публичные
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("NC API\n"))
	})
	mux.HandleFunc("POST /auth/login", authHandler.Login)

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
	mux.Handle("GET /api/events", requireAuth(http.HandlerFunc(eventsHandler.ServeList)))
	mux.Handle("GET /api/events/{id}", requireAuth(http.HandlerFunc(eventsHandler.ServeGet)))
	mux.Handle("GET /api/files", requireAuth(http.HandlerFunc(filesHandler.ServeList)))
	mux.Handle("POST /api/files/upload", requireAuth(http.HandlerFunc(filesHandler.ServeUpload)))
	mux.Handle("DELETE /api/files/{id}", requireAuth(http.HandlerFunc(filesHandler.ServeDelete)))
	mux.Handle("GET /api/chat/ws", requireAuth(http.HandlerFunc(chatWS)))

	// Только преподаватель: оценки, ДЗ
	mux.Handle("POST /api/grades", requireTeacher(http.HandlerFunc(gradesHandler.ServeCreateGrade)))
	mux.Handle("PUT /api/grades/{id}", requireTeacher(http.HandlerFunc(gradesHandler.ServeUpdateGrade)))
	mux.Handle("DELETE /api/grades/{id}", requireTeacher(http.HandlerFunc(gradesHandler.ServeDeleteGrade)))
	mux.Handle("POST /api/schedule/{id}/attachments", requireTeacher(http.HandlerFunc(homeworkHandler.ServeCreate)))
	mux.Handle("PUT /api/attachments/{id}", requireTeacher(http.HandlerFunc(homeworkHandler.ServeUpdate)))
	mux.Handle("DELETE /api/attachments/{id}", requireTeacher(http.HandlerFunc(homeworkHandler.ServeDelete)))

	// Только админ/директор: CRUD событий
	mux.Handle("POST /api/events", requireAdminDirector(http.HandlerFunc(eventsHandler.ServeCreate)))
	mux.Handle("PUT /api/events/{id}", requireAdminDirector(http.HandlerFunc(eventsHandler.ServeUpdate)))
	mux.Handle("DELETE /api/events/{id}", requireAdminDirector(http.HandlerFunc(eventsHandler.ServeDelete)))

	log.Printf("api: listening on %s", cfg.Addr)
	return http.ListenAndServe(cfg.Addr, mux)
}
