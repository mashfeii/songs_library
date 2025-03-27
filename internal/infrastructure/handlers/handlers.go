package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mashfeii/songs_library/internal/application"
	"github.com/mashfeii/songs_library/internal/domain"
	"github.com/sirupsen/logrus"

	clientErrors "github.com/mashfeii/songs_library/internal/infrastructure/errors"
)

// @Summary Get list of songs
// @Description Retrieve list of songs with optional filters and pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Filter by group"
// @Param song query string false "Filter by song"
// @Param song query string false "Filter by text"
// @Param song query string false "Filter by link"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {array} domain.Song "Songs successfully retrieved"
// @Failure 404 {object} domain.ErrorResponse "No songs found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /songs [get]
func GetSongs(service application.SongsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

		if page < 1 {
			logrus.WithFields(logrus.Fields{
				"request_method": c.Request.Method,
				"request_path":   c.Request.URL.Path,
				"request_query":  c.Request.URL.Query(),
			}).Warn("page is less than 1, setting to 1")

			page = 1
		}

		if size < 1 {
			logrus.WithFields(logrus.Fields{
				"request_method": c.Request.Method,
				"request_path":   c.Request.URL.Path,
				"request_query":  c.Request.URL.Query(),
			}).Warn("size is less than 1, setting to 10")

			size = 10
		}

		if _, err := time.Parse("2006-01-02", c.Query("releaseDate")); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":         err,
				"request_query": c.Request.URL.Query(),
			}).Error("Failed to parse release date")

			c.JSON(http.StatusBadRequest, domain.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request",
				Details: "Release date must be in format YYYY-MM-DD",
			})

			return
		}

		filters := map[string]string{
			"group_name":   c.Query("group"),
			"song_name":    c.Query("song"),
			"release_date": c.Query("releaseDate"),
			"text":         c.Query("text"),
			"link":         c.Query("link"),
		}

		songs, err := service.GetSongs(c, filters, page, size)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":          err,
				"request_method": c.Request.Method,
				"request_path":   c.Request.URL.Path,
				"request_query":  c.Request.URL.Query(),
			}).Error("Failed to retrieve songs")

			switch err.(type) {
			case clientErrors.ErrNotFound:
				c.JSON(http.StatusNotFound, domain.ErrorResponse{
					Code:    http.StatusNotFound,
					Message: "Songs not found",
				})
			default:
				c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				})
			}

			return
		}

		logrus.WithFields(logrus.Fields{
			"request_query": c.Request.URL.Query(),
			"amount":        len(songs),
		}).Info("Successfully retrieved songs")
		c.JSON(http.StatusOK, songs)
	}
}

// @Summary Get song verses with pagination
// @Description Retrieve paginated verses of a song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} domain.GetSongVersesResponse "Verses successfully retrieved"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Failure 404 {object} domain.ErrorResponse "No verses found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /songs/{id}/verses [get]
func GetSongVerses(service application.SongsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		size, _ := strconv.Atoi(c.DefaultQuery("size", "1"))

		if page < 1 {
			logrus.WithFields(logrus.Fields{
				"request_method": c.Request.Method,
				"request_path":   c.Request.URL.Path,
				"request_query":  c.Request.URL.Query(),
			}).Warn("page is less than 1, setting to 1")

			page = 1
		}

		if size < 1 {
			logrus.WithFields(logrus.Fields{
				"request_method": c.Request.Method,
				"request_path":   c.Request.URL.Path,
				"request_query":  c.Request.URL.Query(),
			}).Warn("size is less than 1, setting to 10")

			size = 10
		}

		result, err := service.GetSongVerses(c, id, page, size)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":          err,
				"request_method": c.Request.Method,
				"request_path":   c.Request.URL.Path,
				"request_query":  c.Request.URL.Query(),
			}).Error("Failed to retrieve verses")

			switch err.(type) {
			case clientErrors.ErrNotFound:
				c.JSON(http.StatusNotFound, domain.ErrorResponse{
					Code:    http.StatusNotFound,
					Message: "Song not found",
				})
			case clientErrors.ErrInvalidInput:
				c.JSON(http.StatusBadRequest, domain.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: "Page out of range",
					Details: "Start page is greater than the number of verses",
				})
			default:
				c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				})
			}

			return
		}

		logrus.WithFields(logrus.Fields{
			"verses":        result,
			"request_query": c.Request.URL.Query(),
		}).Info("Retrieved verses")
		c.JSON(http.StatusOK, domain.GetSongVersesResponse{
			Verses: result,
			Page:   page,
			Size:   size,
		})
	}
}

// @Summary Delete a song
// @Description Delete a song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 204 "Song successfully removed"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Failure 404 {object} domain.ErrorResponse "Song not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /songs/{id} [delete]
func DeleteSong(service application.SongsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":          err,
				"request_method": c.Request.Method,
				"request_path":   c.Request.URL.Path,
				"request_query":  c.Request.URL.Query(),
			}).Error("Failed to parse song ID")
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request",
				Details: "Song ID must be an integer",
			})
		}

		err = service.DeleteSong(c, id)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":          err,
				"request_method": c.Request.Method,
				"request_path":   c.Request.URL.Path,
				"request_query":  c.Request.URL.Query(),
			}).Error("Failed to remove song")

			switch err.(type) {
			case clientErrors.ErrNotFound:
				c.JSON(http.StatusNotFound, domain.ErrorResponse{
					Code:    http.StatusNotFound,
					Message: "Songs not found",
				})
			default:
				c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				})
			}

			return
		}

		logrus.WithFields(logrus.Fields{
			"id":            id,
			"request_query": c.Request.URL.Query(),
		}).Info("Successfully removed song")
		c.Status(http.StatusNoContent)
	}
}

// @Summary Update a song
// @Description Update a song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body domain.UpdateSongRequest true "Song data"
// @Success 200 {object} domain.Song "Updated song"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Failure 404 {object} domain.ErrorResponse "Song not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /songs/{id} [put]
func UpdateSong(service application.SongsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":          err,
				"request_method": c.Request.Method,
				"request_path":   c.Request.URL.Path,
				"request_query":  c.Request.URL.Query(),
			}).Error("Failed to parse song ID")
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request",
				Details: "Song ID must be an integer",
			})

			return
		}

		var song domain.UpdateSongRequest
		if err := c.BindJSON(&song); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":          err,
				"request_method": c.Request.Method,
				"request_path":   c.Request.URL.Path,
				"request_query":  c.Request.URL.Query(),
			}).Error("Failed to parse song data")
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request",
				Details: "Invalid song data",
			})

			return
		}

		songUpdate := domain.Song{
			ID:          id,
			Song:        song.Song,
			Group:       song.Group,
			ReleaseDate: song.ReleaseDate,
			Text:        song.Text,
			Link:        song.Link,
		}

		err = service.UpdateSong(c, &songUpdate)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":          err,
				"request_method": c.Request.Method,
				"request_path":   c.Request.URL.Path,
				"request_query":  c.Request.URL.Query(),
			}).Error("Failed to update song data")

			switch err.(type) {
			case clientErrors.ErrNotFound:
				c.JSON(http.StatusNotFound, domain.ErrorResponse{
					Code:    http.StatusNotFound,
					Message: "Song not found",
				})
			default:
				c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				})
			}

			return
		}

		logrus.WithFields(logrus.Fields{
			"song":          song,
			"request_query": c.Request.URL.Query(),
		}).Info("Successfully updated song")
		c.JSON(http.StatusOK, song)
	}
}

// @Summary Add a song
// @Description Add a song
// @Tags songs
// @Accept json
// @Produce json
// @Param song body domain.AddSongRequest true "Song name and group"
// @Success 201 {object} domain.AddSongResponse "Song ID"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /songs [post]
func AddSong(service application.SongsServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req domain.AddSongRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		id, err := service.AddSong(c, &req)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":          err,
				"request_method": c.Request.Method,
				"request_path":   c.Request.URL.Path,
				"request_query":  c.Request.URL.Query(),
			}).Error("Failed to add song")

			switch err.(type) {
			case clientErrors.ErrExternal:
				c.JSON(http.StatusNotFound, domain.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "External API error",
				})
			default:
				c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				})
			}

			return
		}

		logrus.WithFields(logrus.Fields{
			"id":            id,
			"request_query": c.Request.URL.Query(),
		}).Info("Successfully added song")
		c.JSON(http.StatusCreated, id)
	}
}
