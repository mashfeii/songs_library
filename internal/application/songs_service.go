package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	client "github.com/mashfeii/songs_library/internal/api"
	"github.com/mashfeii/songs_library/internal/domain"
	clientErrors "github.com/mashfeii/songs_library/internal/infrastructure/errors"
)

type SongsServiceInterface interface {
	GetSongs(ctx context.Context, filters map[string]string, page, size int) ([]domain.Song, error)
	GetSongVerses(ctx context.Context, id, page, size int) ([]string, error)
	DeleteSong(ctx context.Context, id int) error
	UpdateSong(ctx context.Context, song *domain.Song) error
	AddSong(ctx context.Context, songReq *domain.AddSongRequest) (int, error)
}

type SongsService struct {
	repo      domain.SongsRepository
	apiClient client.ClientWithResponsesInterface
}

func NewSongsService(repo domain.SongsRepository, apiClient client.ClientWithResponsesInterface) *SongsService {
	return &SongsService{
		repo:      repo,
		apiClient: apiClient,
	}
}

func (s *SongsService) GetSongs(ctx context.Context, filters map[string]string, page, size int) ([]domain.Song, error) {
	for key, value := range filters {
		if value == "" {
			delete(filters, key)
		}
	}

	songs, err := s.repo.GetSongs(ctx, filters, page, size)
	if err != nil {
		if errors.As(err, &clientErrors.ErrNotFound{}) || len(songs) == 0 {
			return nil, err
		}

		return nil, fmt.Errorf("getting songs: %w", err)
	}

	if len(songs) == 0 {
		return nil, clientErrors.NewErrNotFound("songs")
	}

	return songs, nil
}

func (s *SongsService) GetSongVerses(ctx context.Context, id, page, size int) ([]string, error) {
	song, err := s.repo.GetSongByID(ctx, id)
	if err != nil {
		if errors.As(err, &clientErrors.ErrNotFound{}) {
			return nil, err
		}

		return nil, fmt.Errorf("getting song: %w", err)
	}

	verses := strings.Split(song.Text, "\\n")

	start := (page - 1) * size
	end := start + size

	if start > len(verses) {
		return nil, clientErrors.NewErrInvalidInput("page")
	}

	if end > len(verses) {
		end = len(verses)
	}

	return verses[start:end], nil
}

func (s *SongsService) DeleteSong(ctx context.Context, id int) error {
	return s.repo.DeleteSong(ctx, id)
}

func (s *SongsService) UpdateSong(ctx context.Context, song *domain.Song) error {
	return s.repo.UpdateSong(ctx, song)
}

func (s *SongsService) AddSong(ctx context.Context, songReq *domain.AddSongRequest) (int, error) {
	response, err := s.apiClient.GetInfoWithResponse(ctx,
		&client.GetInfoParams{
			Group: songReq.Group,
			Song:  songReq.Song,
		},
	)
	if err != nil {
		return 0, clientErrors.NewErrExternal(err)
	}

	if response.StatusCode() != http.StatusOK {
		return 0, clientErrors.NewErrExternal(fmt.Errorf("status code: %d", response.StatusCode()))
	}

	song := domain.Song{
		Group:       songReq.Group,
		Song:        songReq.Song,
		ReleaseDate: response.JSON200.ReleaseDate.Time,
		Text:        response.JSON200.Text,
		Link:        response.JSON200.Link,
	}

	id, err := s.repo.AddSong(ctx, &song)
	if err != nil {
		return 0, err
	}

	return id, nil
}
