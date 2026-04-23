package client

import "time"

type EpisodeListing struct {
	Items       []Episode `json:"items"`
	Skip        int       `json:"skip"`
	Take        int       `json:"take"`
	TotalAmount int       `json:"totalAmount"`
}

type Episode struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Program     struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"program"`
	Image struct {
		Square1x1 *ImageSet `json:"square1x1"`
	} `json:"image"`
	Audio struct {
		Broadcast *AudioFile `json:"broadcast"`
		Podcast   *AudioFile `json:"podcast"`
	} `json:"audio"`
	Published time.Time `json:"published"`
	URL       string    `json:"url"`
}

type AudioFile struct {
	ID       int `json:"id"`
	Duration struct {
		Seconds int `json:"seconds"`
	} `json:"duration"`
	MimeType string `json:"mimeType"`
	Variants struct {
		Standard *AudioVariant `json:"standard"`
		High     *AudioVariant `json:"high"`
		Low      *AudioVariant `json:"low"`
	} `json:"variants"`
}

type AudioVariant struct {
	URL     string `json:"url"`
	Bitrate int    `json:"bitrate"`
}

type ImageSet struct {
	Variants []ImageVariant `json:"variants"`
}

type ImageVariant struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type ProgramInfo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       struct {
		Square1x1 *ImageSet `json:"square1x1"`
	} `json:"image"`
	Publisher string `json:"publisher"`
	URL       string `json:"url"`
}
