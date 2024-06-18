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

func convertEpisode(episode Episode) (domain.Episode, error) {
	converted := domain.Episode{
		ID:          episode.ID,
		ProgramID:   episode.Program.ID,
		Title:       episode.Title,
		Description: episode.Description,
		URL:         episode.URL,
		PublishDate: episode.Publishdateutc,
		ImageURL:    episode.Imageurl,
	}

	downloadFile := episode.Downloadpodfile
	broadCastFile := episode.Broadcast.Broadcastfiles.Broadcastfile

	if downloadFile.URL != "" {
		converted.FileURL = downloadFile.URL
		converted.FileDurationSeconds = downloadFile.Duration
		converted.FileBytes = downloadFile.Filesizeinbytes
		converted.ContentType = "audio/mpeg"
	} else if broadCastFile.URL != "" {
		converted.FileURL = broadCastFile.URL
		converted.FileDurationSeconds = broadCastFile.Duration

		contentType, fileSize, err := getFileInfo(broadCastFile.URL)
		if err != nil {
			return converted, errors.WithMessage(err, "could not determine file size")
		}
		converted.FileBytes = fileSize
		converted.ContentType = contentType
	}

	return converted, nil
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
