package api

import (
	"MusicLibrary/internal/repository"
	"MusicLibrary/internal/service"
	"database/sql"
	"github.com/gorilla/mux"
)

// InitializeRoutes - функция инициализации роутера
func InitializeRoutes(db *sql.DB) *mux.Router {
	// Создание репозитория
	repo := repository.NewPostgresSongRepository(db)
	// Создание сервиса
	songService := service.NewSongService(repo)
	// Создание хендлера
	songHandler := NewSongHandler(songService)

	// Создание роутера
	router := mux.NewRouter()

	// Добавление роутеров
	router.HandleFunc("/songs", songHandler.GetSongs).Methods("GET")
	router.HandleFunc("/songs", songHandler.AddSong).Methods("POST")
	router.HandleFunc("/songs/{id}", songHandler.UpdateSong).Methods("PUT")
	router.HandleFunc("/songs/{id}", songHandler.DeleteSong).Methods("DELETE")
	router.HandleFunc("/songs/{id}", songHandler.GetLyrics).Methods("GET")

	return router
}
