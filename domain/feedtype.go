package domain

import "fmt"

// FeedType represents the type of audio file to use for an episode.
type FeedType int

const (
	// FeedTypeDownload represents the shorter download/pod file (without music).
	FeedTypeDownload FeedType = iota
	// FeedTypeBroadcast represents the full broadcast file (with music).
	FeedTypeBroadcast
)

func (f FeedType) String() string {
	switch f {
	case FeedTypeDownload:
		return "download"
	case FeedTypeBroadcast:
		return "broadcast"
	default:
		return fmt.Sprintf("unknown(%d)", int(f))
	}
}
