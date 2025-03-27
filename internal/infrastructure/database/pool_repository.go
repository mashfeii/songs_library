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
)

type PgxRepository struct {
	Pool *pgxpool.Pool
}

func NewPgxRepository(pool *pgxpool.Pool) *PgxRepository {
	return &PgxRepository{Pool: pool}
}

func (r *PgxRepository) GetSongs(ctx context.Context, filters map[string]string, page, size int) ([]domain.Song, error) {
	query := "SELECT id, group_name, song_name, release_date, text, link FROM songs WHERE 1 = 1"
	args := []any{}

	for key, value := range filters {
		query += " AND " + key + " = $" + strconv.Itoa(len(args)+1)
		args = append(args, value)
	}

	query += " LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, size, (page-1)*size)

	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, clientErrors.NewErrNotFound(fmt.Sprintf("songs with args: %v", args))
		}

		return nil, fmt.Errorf("querying songs: %w", clientErrors.NewErrDatabase())
	}
	defer rows.Close()

	var songs []domain.Song

	for rows.Next() {
		var song domain.Song

		err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			return nil, fmt.Errorf("repo scanning songs: %w", clientErrors.NewErrDatabase())
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func (r *PgxRepository) GetSongByID(ctx context.Context, id int) (*domain.Song, error) {
	var song domain.Song

	err := r.Pool.
		QueryRow(ctx, `SELECT id, group_name, song_name, release_date, text, link
    FROM songs
    WHERE id = $1`, id).
		Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, clientErrors.NewErrNotFound(fmt.Sprintf("song with id: %d", id))
		}

		return nil, fmt.Errorf("querying song: %w", clientErrors.NewErrDatabase())
	}

	return &song, nil
}

func (r *PgxRepository) AddSong(ctx context.Context, song *domain.Song) (int, error) {
	var id int

	err := r.Pool.
		QueryRow(ctx, `INSERT INTO songs(group_name, song_name, release_date, text, link)
    VALUES($1, $2, $3, $4, $5)
    RETURNING id`, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link).
		Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repo adding song: %w", clientErrors.NewErrDatabase())
	}

	return id, nil
}

func (r *PgxRepository) UpdateSong(ctx context.Context, song *domain.Song) error {
	_, err := r.Pool.Exec(ctx, `UPDATE songs
    SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5
    WHERE id = $6`, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link, song.ID)

	if errors.Is(err, pgx.ErrNoRows) {
		return clientErrors.NewErrNotFound(fmt.Sprintf("song with id: %d", song.ID))
	}

	return fmt.Errorf("updating song: %w", clientErrors.NewErrDatabase())
}

func (r *PgxRepository) DeleteSong(ctx context.Context, id int) error {
	_, err := r.Pool.Exec(ctx, `DELETE FROM songs WHERE id = $1`, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return clientErrors.NewErrNotFound(fmt.Sprintf("song with id: %d", id))
	}

	if err != nil {
		return fmt.Errorf("deleting song: %w", clientErrors.NewErrDatabase())
	}

	return nil
}
