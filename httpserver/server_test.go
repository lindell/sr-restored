package httpserver

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"context"

	"github.com/lindell/sr-restored/domain"
	"github.com/lindell/sr-restored/podcast"
)

type testRedirectCache struct {
	redirects map[string]string
}

func (c *testRedirectCache) StoreRSS(key string, rawRSS []byte) {}
func (c *testRedirectCache) GetRSS(key string) ([]byte, bool)   { return nil, false }
func (c *testRedirectCache) StoreHash(key string, hash []byte)  {}
func (c *testRedirectCache) GetHash(key string) ([]byte, bool)  { return nil, false }
func (c *testRedirectCache) StoreRedirectURL(key string, redirectURL string) {
	if c.redirects == nil {
		c.redirects = map[string]string{}
	}
	c.redirects[key] = redirectURL
}
func (c *testRedirectCache) GetRedirectURL(key string) (string, bool) {
	v, ok := c.redirects[key]
	return v, ok
}

type testDatabase struct{}

func (db *testDatabase) InsertEpisodes(ctx context.Context, episodes []domain.Episode) error {
	return nil
}
func (db *testDatabase) GetProgram(ctx context.Context, programID int) (domain.Program, error) {
	return domain.Program{}, nil
}
func (db *testDatabase) InsertProgram(ctx context.Context, program domain.Program) error {
	return nil
}
func (db *testDatabase) GetRedirectURL(ctx context.Context, audioFileID int) (string, error) {
	return "", errors.New("not found")
}
func (db *testDatabase) StoreRedirectURL(ctx context.Context, audioFileID int, redirectURL string) error {
	return nil
}

type mockPodcastClient struct {
	resolvedURL string
	err         error
	gotRawURL   string
}

func (m *mockPodcastClient) GetProgram(ctx context.Context, id int, feedTypes []domain.FeedType) (domain.Program, error) {
	return domain.Program{}, errors.New("not implemented")
}

func (m *mockPodcastClient) ResolveRedirectURL(rawURL string) (string, error) {
	m.gotRawURL = rawURL
	if m.err != nil {
		return "", m.err
	}
	return m.resolvedURL, nil
}

func TestRedirectAudioFile(t *testing.T) {
	resolver := &mockPodcastClient{resolvedURL: "https://cdn.example/audio.mp3"}
	s := &Server{Podcast: &podcast.Podcast{Client: resolver, Cache: &testRedirectCache{}, Database: &testDatabase{}}}

	req := httptest.NewRequest(http.MethodGet, "/audio-file/123", nil)
	rec := httptest.NewRecorder()

	s.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMovedPermanently)
	}

	if got := rec.Header().Get("Location"); got != "https://cdn.example/audio.mp3" {
		t.Fatalf("Location = %q, want %q", got, "https://cdn.example/audio.mp3")
	}

	if resolver.gotRawURL != "https://www.sverigesradio.se/topsy/ljudfil/123" {
		t.Fatalf("raw URL = %q, want %q", resolver.gotRawURL, "https://www.sverigesradio.se/topsy/ljudfil/123")
	}
}

func TestRedirectAudioFile_InvalidID(t *testing.T) {
	s := &Server{Podcast: &podcast.Podcast{Client: &mockPodcastClient{}, Cache: &testRedirectCache{}, Database: &testDatabase{}}}

	req := httptest.NewRequest(http.MethodGet, "/audio-file/not-a-number", nil)
	rec := httptest.NewRecorder()

	s.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestRedirectAudioFile_ResolverError(t *testing.T) {
	s := &Server{Podcast: &podcast.Podcast{Client: &mockPodcastClient{err: errors.New("boom")}, Cache: &testRedirectCache{}, Database: &testDatabase{}}}

	req := httptest.NewRequest(http.MethodGet, "/audio-file/123", nil)
	rec := httptest.NewRecorder()

	s.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadGateway)
	}
}
