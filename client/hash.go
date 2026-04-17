package client

import "crypto/sha1"

// Hash generates a hash for the EpisodeListing
func (el EpisodeListing) Hash() []byte {
	hash := sha1.New()

	for _, episode := range el.Items {
		if pod := episode.Audio.Podcast; pod != nil && pod.Variants.Standard != nil {
			hash.Write([]byte(pod.Variants.Standard.URL))
		}
		if bc := episode.Audio.Broadcast; bc != nil && bc.Variants.Standard != nil {
			hash.Write([]byte(bc.Variants.Standard.URL))
		}
	}

	return hash.Sum(nil)
}
