package api

import (
	"MusicLibrary/internal/domain"
	"MusicLibrary/internal/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// SongHandler - структура хендлера для работы с песнями
type SongHandler struct {
	service *service.SongService
}

// NewSongHandler - функция создания хендлера
func NewSongHandler(service *service.SongService) *SongHandler {
	return &SongHandler{service: service}
}

// GetSongs - функция получения списка песен
// @Summary Получение списка песен
// @Description Возвращает список песен с возможностью фильтрации и пагинации
// @Tags songs
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param group query string false "Group"
// @Param song query string false "Song"
// @Param release_date query string false "Release date"
// @Success 200 {array} domain.Song
// @Failure 500 {object} pkg.ErrorResponse
// @Router /songs [get]
func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	// Получение параметров запроса для фильтрации и пагинации
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))
	filter := make(map[string]interface{})

	// Фильтрация по параметрам
	if groupName := query.Get("group"); groupName != "" {
		filter["group_name"] = groupName
	}
	if songName := query.Get("song"); songName != "" {
		filter["song_name"] = songName
	}
	if releaseDate := query.Get("release_date"); releaseDate != "" {
		filter["release_date"] = releaseDate
	}

	// Получение списка песен
	songs, err := h.service.GetSongs(r.Context(), filter, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

// AddSong - функция добавления новой песни
// @Summary Добавление новой песни
// @Description Создает новую песню
// @Tags songs
// @Accept json
// @Produce json
// @Param song body domain.Song true "New song"
// @Success 201
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /songs [post]
func (h *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	// Получение данных из тела запроса
	var newSong domain.Song
	if err := json.NewDecoder(r.Body).Decode(&newSong); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Добавление новой песни
	if err := h.service.AddNewSong(r.Context(), newSong.Group, newSong.Song); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка ответа
	w.WriteHeader(http.StatusCreated)
}

// UpdateSong - функция обновления данных песни
// @Summary Обновление данных песни
// @Description Обновляет информацию о песне по ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body domain.Song true "Updated song"
// @Success 200
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	// Получение данных из тела запроса
	var updatedSong domain.Song
	if err := json.NewDecoder(r.Body).Decode(&updatedSong); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получение ID песни
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}
	updatedSong.ID = id

	// Обновление данных песни
	if err := h.service.UpdateSong(r.Context(), &updatedSong); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка ответа
	w.WriteHeader(http.StatusOK)
}

// DeleteSong - функция удаления песни
// @Summary Удаление песни
// @Description Удаляет песню по ID
// @Tags songs
// @Param id path int true "Song ID"
// @Success 200
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	// Получение ID песни
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	// Удаление песни
	if err := h.service.DeleteSong(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка ответа
	w.WriteHeader(http.StatusOK)
}

// GetLyrics - функция получения текста песни
// @Summary Получение текста песни
// @Description Возвращает текст песни по ID c пагинацией по куплетам
// @Tags songs
// @Param id path int true "Song ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {string} string
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /songs/{id} [get]
func (h *SongHandler) GetLyrics(w http.ResponseWriter, r *http.Request) {
	// Получение ID песни
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	// Получение текста песни
	lyrics, err := h.service.GetSongLyrics(r.Context(), id, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(lyrics))
}
