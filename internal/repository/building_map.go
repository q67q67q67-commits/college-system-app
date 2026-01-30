package repository

import (
	"github.com/q67q67q67-commits/college-system-app/internal/db"
)

// GetBuildingMap возвращает карту здания (первая запись).
func GetBuildingMap() (content, imageURL string, err error) {
	err = db.DB.QueryRow(`SELECT COALESCE(content,''), COALESCE(image_url,'') FROM building_map ORDER BY id LIMIT 1`).Scan(&content, &imageURL)
	return
}

// UpsertBuildingMap создаёт или обновляет карту.
func UpsertBuildingMap(content, imageURL string) error {
	r, err := db.DB.Exec(`UPDATE building_map SET content = $1, image_url = $2, updated_at = now() WHERE id = (SELECT id FROM building_map LIMIT 1)`, content, imageURL)
	if err != nil {
		return err
	}
	n, _ := r.RowsAffected()
	if n == 0 {
		_, err = db.DB.Exec(`INSERT INTO building_map (content, image_url) VALUES ($1, $2)`, content, imageURL)
	}
	return err
}
