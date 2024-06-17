package client

import "github.com/lindell/sr-restored/domain"

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

func convertEpisode(episode Episode) domain.Episode {
	return domain.Episode{
		ID:          episode.ID,
		ProgramID:   episode.Program.ID,
		Title:       episode.Title,
		Description: episode.Description,
		URL:         episode.URL,
		PublishDate: episode.Publishdateutc,
		ImageURL:    episode.Imageurl,

		FileURL:             episode.Downloadpodfile.URL,
		FileDurationSeconds: episode.Downloadpodfile.Duration,
		FileBytes:           episode.Downloadpodfile.Filesizeinbytes,
	}
}
