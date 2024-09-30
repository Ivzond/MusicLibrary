package migrations

import (
	"MusicLibrary/pkg"
	"database/sql"
)

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

-- 	CREATE OR REPLACE FUNCTION update_updated_at_column()
-- 	RETURNS TRIGGER AS $$
-- 	BEGIN
-- 		NEW.updated_at = NOW();
-- 		RETURN NEW;
-- 	END;
-- 	$$ LANGUAGE 'plpgsql';
-- 	
-- 	CREATE TRIGGER update_songs_updated_at 
-- 	BEFORE UPDATE ON songs
-- 	FOR EACH ROW
-- 	EXECUTE PROCEDURE update_updated_at_column();
	`

	_, err := db.Exec(query)
	if err != nil {
		pkg.Error("Ошибка при применении миграций", map[string]interface{}{"error": err})
		return err
	}
	pkg.Info("Миграции успешно применены", nil)
	return nil
}
