package migrations

import (
	"MusicLibrary/pkg"
	"database/sql"
)

// ApplyMigrations - функция применения миграций
func ApplyMigrations(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS songs (
		id SERIAL PRIMARY KEY,
		group_name VARCHAR(255) NOT NULL,
		song_name VARCHAR(255) NOT NULL,
		release_date DATE NOT NULL,
		lyrics TEXT NOT NULL,
		url VARCHAR(500),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_group_name ON songs(group_name);
	CREATE INDEX IF NOT EXISTS idx_song_name ON songs(song_name);
	CREATE INDEX IF NOT EXISTS idx_release_date ON songs(release_date);
`

	_, err := db.Exec(query)
	if err != nil {
		pkg.Error("Ошибка при применении миграций", map[string]interface{}{"error": err})
		return err
	}
	pkg.Info("Миграции успешно применены", nil)
	return nil
}
