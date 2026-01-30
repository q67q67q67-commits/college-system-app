// Narxoz College (NC) — точка входа API-сервера.
//
// Почему код минимальный: это «фундамент» Фазы 1 (Участник Б). Сначала созданы
// структура проекта, схема БД и окружение (docker-compose). Дальше по шагам
// добавляются: сидинг, авторизация (JWT), CRUD-эндпоинты, форум, чаты, файлы.
// Запуск: go run ./cmd/api
package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatalf("api: %v", err)
	}
}

func run(args []string) error {
	// Временный маршрут для проверки: GET / → "NC API"
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("NC API\n"))
	})

	addr := ":8080"
	log.Printf("api: listening on %s", addr)
	return http.ListenAndServe(addr, nil)
}
