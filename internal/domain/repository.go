package domain

import "context"

type SongsRepository interface {
	GetSongs(ctx context.Context, filters map[string]string, page, size int) ([]Song, error)
	GetSongByID(ctx context.Context, id int) (*Song, error)
	AddSong(ctx context.Context, song *Song) (int, error)
	UpdateSong(ctx context.Context, song *Song) error
	DeleteSong(ctx context.Context, id int) error
}
