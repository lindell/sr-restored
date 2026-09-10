package podcast_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/lindell/sr-restored/domain"
	"github.com/lindell/sr-restored/memcache"
	"github.com/lindell/sr-restored/podcast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockClient struct {
	program domain.Program
}

func (m *mockClient) GetProgram(ctx context.Context, id int, feedTypes []domain.FeedType) (domain.Program, error) {
	return m.program, nil
}

type fakeDatabase struct{}

func (f *fakeDatabase) InsertEpisodes(ctx context.Context, episodes []domain.Episode) error {
	return nil
}

func (f *fakeDatabase) GetProgram(ctx context.Context, programID int) (domain.Program, error) {
	return domain.Program{}, nil
}

func (f *fakeDatabase) InsertProgram(ctx context.Context, program domain.Program) error {
	return nil
}

func decompress(t *testing.T, data []byte) string {
	t.Helper()
	gz, err := gzip.NewReader(bytes.NewReader(data))
	require.NoError(t, err)
	defer gz.Close()

	out, err := io.ReadAll(gz)
	require.NoError(t, err)
	return string(out)
}

func TestPodcast_IndentXML(t *testing.T) {
	baseURL, err := url.Parse("https://sr-restored.se/rss")
	require.NoError(t, err)

	program := domain.Program{
		ID:          123,
		Name:        "Test Program",
		Description: "A test podcast program",
		Episodes: []domain.Episode{
			{
				ID:          1,
				Title:       "Episode 1",
				Description: "First episode",
			},
		},
	}

	t.Run("IndentXML false (production mode)", func(t *testing.T) {
		p := &podcast.Podcast{
			Client:    &mockClient{program: program},
			Cache:     memcache.NewCache(),
			Database:  &fakeDatabase{},
			RSSUrl:    baseURL,
			Now:       func() time.Time { return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) },
			IndentXML: false,
		}

		raw, _, err := p.GetPodcast(context.Background(), 123, nil)
		require.NoError(t, err)

		xmlStr := decompress(t, raw)

		// In production mode, the XML is marshaled without newlines between elements
		assert.False(t, strings.Contains(xmlStr, "\n"), "production XML should be a single line without newlines")
	})

	t.Run("IndentXML true (test mode)", func(t *testing.T) {
		p := &podcast.Podcast{
			Client:    &mockClient{program: program},
			Cache:     memcache.NewCache(),
			Database:  &fakeDatabase{},
			RSSUrl:    baseURL,
			Now:       func() time.Time { return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) },
			IndentXML: true,
		}

		raw, _, err := p.GetPodcast(context.Background(), 123, nil)
		require.NoError(t, err)

		xmlStr := decompress(t, raw)

		// In test mode, the XML is marshaled with newlines and indentation
		assert.True(t, strings.Contains(xmlStr, "\n"), "test XML should contain newlines")
		assert.True(t, strings.Contains(xmlStr, "\n  <channel>"), "test XML should have indented channel element")
	})
}
