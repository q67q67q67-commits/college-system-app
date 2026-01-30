package config

import (
	"os"
)

// Config — настройки приложения из переменных окружения.
type Config struct {
	Addr       string // HTTP-адрес (например :8080)
	DBURL      string // PostgreSQL connection string
	JWTSecret  string // Секрет для подписи JWT
	UploadPath string // Папка для загруженных файлов
}

// Load читает конфиг из окружения.
func Load() *Config {
	addr := os.Getenv("NC_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	dbURL := os.Getenv("NC_DB_URL")
	if dbURL == "" {
		dbURL = "postgres://nc:nc_dev@localhost:5432/nc_db?sslmode=disable"
	}
	jwtSecret := os.Getenv("NC_JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "nc-dev-secret-change-in-production"
	}
	uploadPath := os.Getenv("NC_UPLOAD_PATH")
	if uploadPath == "" {
		uploadPath = "./uploads"
	}
	return &Config{
		Addr:       addr,
		DBURL:      dbURL,
		JWTSecret:  jwtSecret,
		UploadPath: uploadPath,
	}
}
