package repository

import (
	"MusicLibrary/internal/domain"
	"context"
)

// SongRepository - интерфейс репозитория для реализации логики взаимодействияс БД
type SongRepository interface {
	CreateSong(ctx context.Context, song *domain.Song) error
	GetSongs(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]domain.Song, error)
	GetSongByID(ctx context.Context, id int) (*domain.Song, error)
	UpdateSong(ctx context.Context, song *domain.Song) error
	DeleteSong(ctx context.Context, id int) error
	GetSongLyrics(ctx context.Context, id, limit, offset int) (string, error)
}
