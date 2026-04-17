package client

import (
	"net/http"
	"strconv"

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
			contentType, size, err := c.getFileInfo(pod.Variants.Standard.URL)
			if err != nil {
				return fileInfo{}, errors.WithMessage(err, "could not determine podcast file info")
			}
			return fileInfo{
				URL:             pod.Variants.Standard.URL,
				DurationSeconds: pod.Duration.Seconds,
				Bytes:           size,
				ContentType:     contentType,
			}, nil

		case domain.FeedTypeBroadcast:
			bc := episode.Audio.Broadcast
			if bc == nil || bc.Variants.Standard == nil || bc.Variants.Standard.URL == "" {
				continue
			}
			contentType, size, err := c.getFileInfo(bc.Variants.Standard.URL)
			if err != nil {
				return fileInfo{}, errors.WithMessage(err, "could not determine broadcast file info")
			}
			return fileInfo{
				URL:             bc.Variants.Standard.URL,
				DurationSeconds: bc.Duration.Seconds,
				Bytes:           size,
				ContentType:     contentType,
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

func (c *Client) getFileInfo(fileURL string) (contentType string, size int, err error) {
	if contentType, size, ok := c.Cache.GetFileInfo(fileURL); ok {
		return contentType, size, nil
	}

	req, err := http.NewRequest(http.MethodHead, fileURL, nil)
	if err != nil {
		return "", 0, errors.WithMessage(err, "could not create HEAD request")
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", 0, errors.WithMessage(err, "could not fetch file url")
	}

	contentLength := res.Header.Get("Content-Length")
	size, err = strconv.Atoi(contentLength)
	if err != nil {
		return "", 0, errors.WithMessage(err, "could not parse file url content length")
	}

	contentType = res.Header.Get("Content-Type")

	c.Cache.StoreFileInfo(fileURL, contentType, size)

	return contentType, size, nil
}
