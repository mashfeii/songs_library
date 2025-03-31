package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mashfeii/songs_library/internal/application"
	"github.com/mashfeii/songs_library/internal/domain"
	"github.com/mashfeii/songs_library/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	clientErrors "github.com/mashfeii/songs_library/internal/infrastructure/errors"
)

func TestSongsService_GetSongs(t *testing.T) {
	mockSongsRepo := mocks.NewSongsRepositoryMock(t)
	mockGroupsRepo := mocks.NewGroupsRepositoryMock(t)

	service := application.NewSongsService(mockSongsRepo, mockGroupsRepo, nil)

	t.Run("Success", func(t *testing.T) {
		mockSongsRepo.On("GetSongs", mock.Anything, mock.Anything, 1, 10).Return([]domain.Song{
			{
				ID:          1,
				GroupID:     1,
				Group:       "Muse",
				Song:        "Supermassive Black Hole",
				ReleaseDate: time.Date(2006, 6, 19, 0, 0, 0, 0, time.UTC),
				Text:        "Oh baby dont you know I suffer",
				Link:        "https://www.youtube.com/watch?v=UqLRqzTp6Rk",
			},
		}, nil).Once()

		songs, err := service.GetSongs(context.Background(), map[string]string{}, 1, 10)
		assert.NoError(t, err)
		assert.Len(t, songs, 1)
		assert.Equal(t, "Muse", songs[0].Group)
		mockSongsRepo.AssertExpectations(t)
	})

	t.Run("NoSongsFound", func(t *testing.T) {
		mockSongsRepo.On("GetSongs", mock.Anything, mock.Anything, 1, 10).Return([]domain.Song{}, nil).Once()

		songs, err := service.GetSongs(context.Background(), map[string]string{}, 1, 10)
		assert.Error(t, err)
		assert.Nil(t, songs)
		assert.True(t, errors.As(err, &clientErrors.ErrNotFound{}))
		mockSongsRepo.AssertExpectations(t)
	})
}
