package api

import (
	"MusicLibrary/internal/repository"
	"MusicLibrary/internal/service"
	"database/sql"
	"github.com/gorilla/mux"
)

func InitializeRoutes(db *sql.DB) *mux.Router {
	repo := repository.NewPostgresSongRepository(db)
	songService := service.NewSongService(repo)
	songHandler := NewSongHandler(songService)

	router := mux.NewRouter()

	router.HandleFunc("/songs", songHandler.GetSongs).Methods("GET")
	router.HandleFunc("/songs", songHandler.AddSong).Methods("POST")
	router.HandleFunc("/songs/{id}", songHandler.UpdateSong).Methods("PUT")
	router.HandleFunc("/songs/{id}", songHandler.DeleteSong).Methods("DELETE")
	router.HandleFunc("/songs/{id}", songHandler.GetLyrics).Methods("GET")

	return router
}
