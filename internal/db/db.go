package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DB — глобальный пул подключений (для простоты; в продакшене — через dependency injection).
var DB *sql.DB

// Open подключается к PostgreSQL по connURL и сохраняет в DB.
func Open(connURL string) error {
	var err error
	DB, err = sql.Open("postgres", connURL)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("ping: %w", err)
	}
	return nil
}

// Close закрывает пул.
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
