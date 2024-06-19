package podcast

import (
	"slices"

	"github.com/lindell/sr-restored/domain"
)

func mergeEpisodes(eps1, eps2 []domain.Episode) []domain.Episode {
	episodeMap := map[int]domain.Episode{}

	for _, episode := range eps1 {
		episodeMap[episode.ID] = episode
	}

	for _, episode := range eps2 {
		episodeMap[episode.ID] = episode
	}

	episodes := make([]domain.Episode, 0, len(episodeMap))
	for _, episode := range episodeMap {
		episodes = append(episodes, episode)
	}

	slices.SortFunc(episodes, func(a, b domain.Episode) int {
		return b.PublishDate.Compare(a.PublishDate)
	})

	return episodes
}
