package service

import (
	"MusicLibrary/internal/domain"
	"MusicLibrary/internal/repository"
	"MusicLibrary/pkg"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type SongService struct {
	repo repository.SongRepository
}

func NewSongService(repo repository.SongRepository) *SongService {
	return &SongService{repo: repo}
}

func (s *SongService) GetSongs(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]domain.Song, error) {
	pkg.Info("Получение списка песен с фильтрацией", filters)
	songs, err := s.repo.GetSongs(ctx, filters, limit, offset)
	if err != nil {
		pkg.Error("Ошибка при получении песен", map[string]interface{}{"error": err})
		return nil, err
	}
	pkg.Debug(fmt.Sprintf("Найдено песен: %d", len(songs)), nil)
	return songs, nil
}

func (s *SongService) GetSongLyrics(ctx context.Context, id, limit, offset int) (string, error) {
	pkg.Info("Получение текста песни", map[string]interface{}{"id": id, "limit": limit, "offset": offset})
	lyrics, err := s.repo.GetSongLyrics(ctx, id, limit, offset)
	if err != nil {
		pkg.Error("Ошибка при получении текста песни", map[string]interface{}{"error": err})
		return "", err
	}
	pkg.Debug("Получен текст песни", map[string]interface{}{"lyrics_length": len(lyrics)})
	return lyrics, err
}

func (s *SongService) DeleteSong(ctx context.Context, id int) error {
	pkg.Info("Удаление песни", map[string]interface{}{"id": id})
	err := s.repo.DeleteSong(ctx, id)
	if err != nil {
		pkg.Error("Ошибка при удалении песни", map[string]interface{}{"error": err})
		return err
	}
	pkg.Debug("Песня успешно удалена", nil)
	return nil
}

func (s *SongService) UpdateSong(ctx context.Context, updatedSong *domain.Song) error {
	pkg.Info("Обновление данных песни", map[string]interface{}{"song_id": updatedSong.ID})

	currentSong, err := s.repo.GetSongByID(ctx, updatedSong.ID)
	fmt.Println(currentSong)
	if err != nil {
		pkg.Error("Ошибка при получении текущей версии песни", map[string]interface{}{"error": err})
		return err
	}

	if updatedSong.Group != "" {
		currentSong.Group = updatedSong.Group
	}
	if updatedSong.Song != "" {
		currentSong.Song = updatedSong.Song
	}
	if updatedSong.Lyrics != "" {
		currentSong.Lyrics = updatedSong.Lyrics
	}
	if updatedSong.ReleaseDate.IsZero() {
		updatedSong.ReleaseDate = currentSong.ReleaseDate
	}
	if updatedSong.URL != "" {
		currentSong.URL = updatedSong.URL
	}

	err = s.repo.UpdateSong(ctx, currentSong)
	if err != nil {
		pkg.Error("Ошибка при обновлении песни", map[string]interface{}{"error": err})
		return err
	}
	pkg.Debug("Данные успешно обновлены", map[string]interface{}{"song_id": currentSong.ID})
	return nil
}

func (s *SongService) AddNewSong(ctx context.Context, group, song string) error {
	pkg.Info("Добавление новой песни", map[string]interface{}{"group": group, "song": song})

	// Получение информации о песне через внешний API
	info, err := s.fetchSongInfo(group, song)
	if err != nil {
		pkg.Error("Ошибка при запросе данных из внешнего API", map[string]interface{}{"error": err})
		return err
	}

	// Создание новой сущности песни на основе данных API
	newSong := &domain.Song{
		Group:       group,
		Song:        song,
		ReleaseDate: info.ReleaseDate,
		Lyrics:      info.Text,
		URL:         info.Link,
	}

	// Сохранение в базу данных
	err = s.repo.CreateSong(ctx, newSong)
	if err != nil {
		pkg.Error("Ошибка при сохранении песни в базу данных", map[string]interface{}{"error": err})
		return err
	}

	pkg.Debug("Новая песня успешно добавлена", map[string]interface{}{"song_id": newSong.ID})
	return nil
}

func (s *SongService) fetchSongInfo(group, song string) (*domain.SongDetail, error) {
	pkg.Info("Запрос данных о песне через внешний API", map[string]interface{}{"group": group, "song": song})

	// Формируем URL запроса к API
	url := fmt.Sprintf("http://%s/info?group=%s&song=%s", os.Getenv("EXTERNAL_API_URL"), group, song)

	resp, err := http.Get(url)
	if err != nil {
		pkg.Error("Ошибка HTTP запроса", map[string]interface{}{"error": err})
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		pkg.Error("Внешнее API вернуло ошибку", map[string]interface{}{"status_code": resp.StatusCode})
		return nil, fmt.Errorf("не удалось запросить данные о песне, статус-код: %d", resp.StatusCode)
	}

	var info domain.SongDetail
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		pkg.Error("Ошибка декодирования ответа API", map[string]interface{}{"error": err})
		return nil, err
	}

	pkg.Debug("Данные о песне успешно получены", map[string]interface{}{"group": group, "song": song})
	return &info, nil
}
