package database

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mashfeii/songs_library/internal/domain"
	clientErrors "github.com/mashfeii/songs_library/internal/infrastructure/errors"
	"github.com/sirupsen/logrus"
)

type SongsRepository interface {
	GetSongs(ctx context.Context, filters map[string]string, page, size int) ([]domain.Song, error)
	GetSongByID(ctx context.Context, id int) (*domain.Song, error)
	AddSong(ctx context.Context, song *domain.Song) (int, error)
	UpdateSong(ctx context.Context, song *domain.Song) error
	DeleteSong(ctx context.Context, id int) error
}

type SongsPoolRepository struct {
	Pool *pgxpool.Pool
}

func NewSongsPoolRepository(pool *pgxpool.Pool) *SongsPoolRepository {
	return &SongsPoolRepository{Pool: pool}
}

func (r *SongsPoolRepository) GetSongs(ctx context.Context, filters map[string]string, page, size int) ([]domain.Song, error) {
	query := `
    SELECT s.id, group_id, g.name as group_name, song_name, release_date, text, link
    FROM songs AS s
    JOIN groups AS g ON s.group_id = g.id
    WHERE 1 = 1`
	args := []any{}

	for key, value := range filters {
		query += " AND " + key + " = $" + strconv.Itoa(len(args)+1)
		args = append(args, value)
	}

	query += " LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, size, (page-1)*size)

	logrus.WithFields(logrus.Fields{
		"query": query,
		"args":  args,
	}).Debug("Executing get songs query")

	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"filters": filters,
		}).Error("Failed to get songs from database")

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, clientErrors.NewErrNotFound(fmt.Sprintf("songs with args: %v", args))
		}

		return nil, fmt.Errorf("querying songs: %w", clientErrors.NewErrDatabase())
	}
	defer rows.Close()

	var songs []domain.Song

	for rows.Next() {
		var song domain.Song

		err := rows.Scan(&song.ID, &song.GroupID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			return nil, fmt.Errorf("repo scanning songs: %w", clientErrors.NewErrDatabase())
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func (r *SongsPoolRepository) GetSongByID(ctx context.Context, id int) (*domain.Song, error) {
	logrus.WithFields(logrus.Fields{
		"id": id,
	}).Debug("Executing get song by id query")

	var song domain.Song

	err := r.Pool.
		QueryRow(ctx, `
    SELECT s.id, group_id, g.name AS group_name, song_name, release_date, text, link
    FROM songs AS s 
    JOIN groups AS g ON s.group_id = g.id
    WHERE s.id = $1
    `, id).
		Scan(&song.ID, &song.GroupID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"id":    id,
		}).Error("Failed to get song from database")

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, clientErrors.NewErrNotFound(fmt.Sprintf("song with id: %d", id))
		}

		return nil, fmt.Errorf("querying song: %w", clientErrors.NewErrDatabase())
	}

	return &song, nil
}

func (r *SongsPoolRepository) AddSong(ctx context.Context, song *domain.Song) (int, error) {
	logrus.WithFields(logrus.Fields{
		"song": song,
	}).Debug("Executing add song query")

	var id int

	err := r.Pool.
		QueryRow(ctx, `INSERT INTO songs(group_id, song_name, release_date, text, link)
    VALUES($1, $2, $3, $4, $5)
    RETURNING id`, song.GroupID, song.Song, song.ReleaseDate, song.Text, song.Link).
		Scan(&id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"song":  song,
		}).Error("Failed to add song to database")

		return 0, fmt.Errorf("repo adding song: %w", clientErrors.NewErrDatabase())
	}

	return id, nil
}

func (r *SongsPoolRepository) UpdateSong(ctx context.Context, song *domain.Song) error {
	logrus.WithFields(logrus.Fields{
		"song": song,
	}).Debug("Executing update song query")

	_, err := r.Pool.Exec(ctx, `UPDATE songs
    SET group_id = $1, song_name = $2, release_date = $3, text = $4, link = $5
    WHERE id = $6`, song.GroupID, song.Song, song.ReleaseDate, song.Text, song.Link, song.ID)

	if errors.Is(err, pgx.ErrNoRows) {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"song":  song,
		}).Error("Failed to update song in database")

		return clientErrors.NewErrNotFound(fmt.Sprintf("song with id: %d", song.ID))
	}

	return fmt.Errorf("updating song: %w", clientErrors.NewErrDatabase())
}

func (r *SongsPoolRepository) DeleteSong(ctx context.Context, id int) error {
	logrus.WithFields(logrus.Fields{
		"id": id,
	}).Debug("Executing delete song query")

	_, err := r.Pool.Exec(ctx, `DELETE FROM songs WHERE id = $1`, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return clientErrors.NewErrNotFound(fmt.Sprintf("song with id: %d", id))
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"id":    id,
		}).Error("Failed to delete song from database")

		return fmt.Errorf("deleting song: %w", clientErrors.NewErrDatabase())
	}

	return nil
}
