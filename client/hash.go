package client

import "crypto/sha1"

// Hash generates a hash for the EpisodeListing
func (el EpisodeListing) Hash() []byte {
	hash := sha1.New()

	for _, episode := range el.Episodes.Episode {
		// If all files are the same, we assume the listing is the same
		hash.Write([]byte(episode.Listenpodfile.URL))
		hash.Write([]byte(episode.Downloadpodfile.URL))
		hash.Write([]byte(episode.Broadcast.Broadcastfiles.Broadcastfile.URL))
	}

	return hash.Sum(nil)
}
