package client

import (
	"github.com/lindell/sr-restored/domain"
	"github.com/pkg/errors"
)

func convertProgram(program ProgramInfo) domain.Program {
	var imageURL string
	if program.Image.Square1x1 != nil && len(program.Image.Square1x1.Variants) > 0 {
		imageURL = program.Image.Square1x1.Variants[0].URL
	}

	return domain.Program{
		ID:          program.ID,
		Name:        program.Title,
		Description: program.Description,
		Copyright:   "Copyright Sveriges Radio. All rights reserved.",
		URL:         program.URL,
		ImageURL:    imageURL,
	}
}

// fileInfo holds the resolved audio file metadata for an episode.
type fileInfo struct {
	URL             string
	DurationSeconds int
	Bytes           int
	ContentType     string
}

func (c *Client) convertEpisode(episode Episode, feedTypes []domain.FeedType) (domain.Episode, error) {
	var imageURL string
	if episode.Image.Square1x1 != nil && len(episode.Image.Square1x1.Variants) > 0 {
		imageURL = episode.Image.Square1x1.Variants[0].URL
	}

	converted := domain.Episode{
		ID:          episode.ID,
		ProgramID:   episode.Program.ID,
		Title:       episode.Title,
		Description: episode.Description,
		URL:         episode.URL,
		PublishDate: episode.Published,
		ImageURL:    imageURL,
	}

	fi, err := c.resolveFileInfo(episode, feedTypes)
	if err != nil {
		return converted, err
	}

	converted.FileURL = fi.URL
	converted.FileDurationSeconds = fi.DurationSeconds
	converted.FileBytes = fi.Bytes
	converted.ContentType = fi.ContentType

	return converted, nil
}

// resolveFileInfo selects the best available audio file for an episode based on
// the preferred feed types. It returns nil if no suitable file is found.
func (c *Client) resolveFileInfo(episode Episode, feedTypes []domain.FeedType) (fileInfo, error) {
	for _, ft := range feedTypes {
		switch ft {
		case domain.FeedTypeDownload:
			pod := episode.Audio.Podcast
			if pod == nil || pod.Variants.Standard == nil || pod.Variants.Standard.URL == "" {
				continue
			}
			return fileInfo{
				URL:             pod.Variants.Standard.URL,
				DurationSeconds: pod.Duration.Seconds,
				Bytes:           estimateFileSize(pod.Variants.Standard.Bitrate, pod.Duration.Seconds),
				ContentType:     pod.MimeType,
			}, nil

		case domain.FeedTypeBroadcast:
			bc := episode.Audio.Broadcast
			if bc == nil || bc.Variants.Standard == nil || bc.Variants.Standard.URL == "" {
				continue
			}
			return fileInfo{
				URL:             bc.Variants.Standard.URL,
				DurationSeconds: bc.Duration.Seconds,
				Bytes:           estimateFileSize(bc.Variants.Standard.Bitrate, bc.Duration.Seconds),
				ContentType:     bc.MimeType,
			}, nil
		}
	}

	return fileInfo{}, errors.New("could not find download file")
}

// FileInfoCache is an optional cache for file metadata lookups.
type FileInfoCache interface {
	StoreFileInfo(key string, contentType string, size int)
	GetFileInfo(key string) (contentType string, size int, ok bool)
}

func estimateFileSize(bitrate int, durationSeconds int) int {
	if bitrate <= 0 || durationSeconds <= 0 {
		return 0
	}

	return bitrate * durationSeconds / 8
}
