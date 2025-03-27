package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mashfeii/songs_library/internal/domain"
	"github.com/mashfeii/songs_library/internal/infrastructure/handlers"
	"github.com/mashfeii/songs_library/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	clientErrors "github.com/mashfeii/songs_library/internal/infrastructure/errors"
)

func TestGetSongs(t *testing.T) {
	mockService := mocks.NewSongsServiceInterfaceMock(t)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/songs?page=1&size=10", http.NoBody)

	t.Run("Success", func(t *testing.T) {
		mockService.On("GetSongs", mock.Anything, mock.Anything, 1, 10).Return([]domain.Song{
			{ID: 1, Group: "Muse", Song: "Supermassive Black Hole"},
		}, nil).Once()

		handlers.GetSongs(mockService)(c)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[{"id":1,"group":"Muse","song":"Supermassive Black Hole","release_date":"0001-01-01T00:00:00Z","text":"","link":""}]`, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/songs?page=1&size=10", http.NoBody)

		mockService.On("GetSongs", mock.Anything, mock.Anything, 1, 10).Return(nil, clientErrors.NewErrNotFound("songs")).Once()

		handlers.GetSongs(mockService)(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"code":404,"message":"Songs not found"}`, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/songs?page=1&size=10", http.NoBody)

		mockService.On("GetSongs", mock.Anything, mock.Anything, 1, 10).Return(nil, clientErrors.NewErrDatabase()).Once()

		handlers.GetSongs(mockService)(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"code":500,"message":"Internal server error"}`, w.Body.String())
		mockService.AssertExpectations(t)
	})
}
