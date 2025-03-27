package domain

import "time"

type AddSongRequest struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type UpdateSongRequest struct {
	Group       string    `json:"group"`
	Song        string    `json:"song"`
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}
