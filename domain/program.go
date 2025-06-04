package domain

type Program struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Email       string `db:"email"`
	Copyright   string `db:"copyright"`
	URL         string `db:"url"`
	ImageURL    string `db:"image_url"`

	Hash []byte

	Episodes []Episode
}
