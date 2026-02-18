package client

import (
	"net/http"
	"strconv"

	"github.com/lindell/sr-restored/domain"
	"github.com/pkg/errors"
)

func convertProgram(program ProgramInfo) domain.Program {
	return domain.Program{
		ID:          program.Program.ID,
		Name:        program.Program.Name,
		Description: program.Program.Description,
		Email:       program.Program.Email,
		Copyright:   program.Copyright,
		URL:         program.Program.Programurl,
		ImageURL:    program.Program.Programimage,
	}
}

// fileInfo holds the resolved audio file metadata for an episode.
type fileInfo struct {
	URL             string
	DurationSeconds int
	Bytes           int
	ContentType     string
}

func convertEpisode(episode Episode, feedTypes []domain.FeedType) (domain.Episode, error) {
	converted := domain.Episode{
		ID:          episode.ID,
		ProgramID:   episode.Program.ID,
		Title:       episode.Title,
		Description: episode.Description,
		URL:         episode.URL,
		PublishDate: episode.Publishdateutc,
		ImageURL:    episode.Imageurl,
	}

	fi, err := resolveFileInfo(episode, feedTypes)
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
func resolveFileInfo(episode Episode, feedTypes []domain.FeedType) (fileInfo, error) {
	for _, ft := range feedTypes {
		switch ft {
		case domain.FeedTypeDownload:
			dl := episode.Downloadpodfile
			if dl.URL == "" {
				continue
			}
			return fileInfo{
				URL:             dl.URL,
				DurationSeconds: dl.Duration,
				Bytes:           dl.Filesizeinbytes,
				ContentType:     "audio/mpeg",
			}, nil

		case domain.FeedTypeBroadcast:
			bc := episode.Broadcast.Broadcastfiles.Broadcastfile
			if bc.URL == "" {
				continue
			}
			contentType, size, err := getFileInfo(bc.URL)
			if err != nil {
				return fileInfo{}, errors.WithMessage(err, "could not determine broadcast file info")
			}
			return fileInfo{
				URL:             bc.URL,
				DurationSeconds: bc.Duration,
				Bytes:           size,
				ContentType:     contentType,
			}, nil
		}
	}

	return fileInfo{}, errors.New("could not find download file")
}

func getFileInfo(fileURL string) (contentType string, size int, err error) {
	res, err := http.Head(fileURL)
	if err != nil {
		return "", 0, errors.WithMessage(err, "could not fetch file url")
	}

	contentLength := res.Header.Get("Content-Length")
	size, err = strconv.Atoi(contentLength)
	if err != nil {
		return "", 0, errors.WithMessage(err, "could not parse file url content length")
	}

	return res.Header.Get("Content-Type"), size, nil
}
