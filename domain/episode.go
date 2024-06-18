package domain

import "time"

type Episode struct {
	ID          int       `db:"id"`
	ProgramID   int       `db:"program_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	URL         string    `db:"url"`
	PublishDate time.Time `db:"publish_date"`
	ImageURL    string    `db:"image_url"`
	ContentType string

	FileURL             string `db:"file_url"`
	FileDurationSeconds int    `db:"file_duration"`
	FileBytes           int    `db:"file_bytes"`
}
