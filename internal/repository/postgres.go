package repository

import (
	"MusicLibrary/internal/domain"
	"MusicLibrary/pkg"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type PostgresSongRepository struct {
	db *sql.DB
}

func NewPostgresSongRepository(db *sql.DB) *PostgresSongRepository {
	return &PostgresSongRepository{db: db}
}

func (r *PostgresSongRepository) CreateSong(ctx context.Context, song *domain.Song) error {
	query := `INSERT INTO songs (group_name, song_name, release_date, lyrics, url, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRowContext(ctx, query, song.Group, song.Song, song.ReleaseDate.Format("2006-01-02"), song.Lyrics, song.URL, time.Now(), time.Now()).Scan(&song.ID)
}

func (r *PostgresSongRepository) GetSongs(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]domain.Song, error) {
	query := `SELECT id, group_name, song_name, release_date, lyrics, url, created_at, updated_at FROM songs WHERE 1=1`
	var args []interface{}
	i := 1

	// Динамически добавляем фильтры
	for key, value := range filter {
		query += fmt.Sprintf(" AND %s = $%d", key, i)
		args = append(args, value)
		i++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, limit, offset)

	fmt.Println("Executing query:", query, args)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		pkg.Error("Ошибка при выполнении запроса", map[string]interface{}{"error": err})
		return nil, err
	}
	defer rows.Close()

	var songs []domain.Song
	for rows.Next() {
		var song domain.Song
		err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Lyrics, &song.URL, &song.CreatedAt, &song.UpdatedAt)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func (r *PostgresSongRepository) GetSongByID(ctx context.Context, id int) (*domain.Song, error) {
	query := `SELECT id, group_name, song_name, release_date, lyrics, url, created_at, updated_at FROM songs WHERE id = $1`
	var song domain.Song
	err := r.db.QueryRowContext(ctx, query, id).Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Lyrics, &song.URL, &song.CreatedAt, &song.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &song, nil
}

func (r *PostgresSongRepository) UpdateSong(ctx context.Context, song *domain.Song) error {
	query := `UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, lyrics = $4, url = $5, updated_at = $6 WHERE id = $7`
	_, err := r.db.ExecContext(ctx, query, song.Group, song.Song, song.ReleaseDate.Format("2006-01-02"), song.Lyrics, song.URL, time.Now(), song.ID)
	return err
}

func (r *PostgresSongRepository) DeleteSong(ctx context.Context, id int) error {
	query := `DELETE FROM songs WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresSongRepository) GetSongLyrics(ctx context.Context, id, limit, offset int) (string, error) {
	query := `SELECT lyrics FROM songs WHERE id = $1`
	var lyrics string
	err := r.db.QueryRowContext(ctx, query, id).Scan(&lyrics)
	if err != nil {
		return "", err
	}

	verses := splitLyrics(lyrics)
	start := offset
	end := offset + limit
	if end > len(verses) {
		end = len(verses)
	}
	return combineLyrics(verses[start:end]), nil
}

func splitLyrics(lyrics string) []string {
	return strings.Split(lyrics, "\n\n")
}

func combineLyrics(verses []string) string {
	return strings.Join(verses, "\n\n")
}
