package podcast

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/lindell/go-stderrs/stderrs"
	"github.com/lindell/sr-restored/domain"
)

func TestToAudioProxyURL(t *testing.T) {
	rssURL, err := url.Parse("https://sr-restored.se/rss")
	if err != nil {
		t.Fatal(err)
	}

	p := &Podcast{RSSUrl: rssURL}

	tests := []struct {
		name   string
		rawURL string
		want   string
	}{
		{
			name:   "rewrites topsy url",
			rawURL: "https://www.sverigesradio.se/topsy/ljudfil/42",
			want:   "https://sr-restored.se/audio-file/42",
		},
		{
			name:   "keeps non topsy url",
			rawURL: "https://static-cdn.sr.se/laddahem/podradio/test.mp3",
			want:   "https://static-cdn.sr.se/laddahem/podradio/test.mp3",
		},
		{
			name:   "keeps invalid topsy id",
			rawURL: "https://www.sverigesradio.se/topsy/ljudfil/abc",
			want:   "https://www.sverigesradio.se/topsy/ljudfil/abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.toAudioProxyURL(tt.rawURL)
			if got != tt.want {
				t.Fatalf("toAudioProxyURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

type testRedirectCache struct {
	mu        sync.Mutex
	redirects map[string]string
}

func (c *testRedirectCache) StoreRSS(key string, rawRSS []byte) {}
func (c *testRedirectCache) GetRSS(key string) ([]byte, bool)   { return nil, false }
func (c *testRedirectCache) StoreHash(key string, hash []byte)  {}
func (c *testRedirectCache) GetHash(key string) ([]byte, bool)  { return nil, false }
func (c *testRedirectCache) StoreRedirectURL(key string, redirectURL string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.redirects == nil {
		c.redirects = map[string]string{}
	}
	c.redirects[key] = redirectURL
}
func (c *testRedirectCache) GetRedirectURL(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.redirects == nil {
		return "", false
	}
	v, ok := c.redirects[key]
	return v, ok
}

type testClient struct {
	resolveFunc func(rawURL string) (string, error)
}

type testDatabase struct {
	mu               sync.Mutex
	redirects        map[int]string
	storedRedirects  map[int]string
	getRedirectCalls int
	storeCalls       int
}

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
	db.mu.Lock()
	defer db.mu.Unlock()
	db.getRedirectCalls++
	if redirectURL, ok := db.redirects[audioFileID]; ok {
		return redirectURL, nil
	}
	return "", stderrs.NewNotFound("not found")
}

func (db *testDatabase) StoreRedirectURL(ctx context.Context, audioFileID int, redirectURL string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.storeCalls++
	if db.storedRedirects == nil {
		db.storedRedirects = map[int]string{}
	}
	db.storedRedirects[audioFileID] = redirectURL
	return nil
}

func (c *testClient) GetProgram(ctx context.Context, id int, feedTypes []domain.FeedType) (domain.Program, error) {
	return domain.Program{}, nil
}

func (c *testClient) ResolveRedirectURL(rawURL string) (string, error) {
	if c.resolveFunc != nil {
		return c.resolveFunc(rawURL)
	}
	return rawURL, nil
}

func TestResolveRedirectURL_CachesLookup(t *testing.T) {
	cache := &testRedirectCache{}
	var calls atomic.Int32

	p := &Podcast{
		Cache:    cache,
		Database: &testDatabase{},
		Client: &testClient{resolveFunc: func(rawURL string) (string, error) {
			calls.Add(1)
			return "https://cdn.example/audio-1.mp3", nil
		}},
	}

	first, err := p.ResolveRedirectURL(1)
	if err != nil {
		t.Fatal(err)
	}
	second, err := p.ResolveRedirectURL(1)
	if err != nil {
		t.Fatal(err)
	}

	if first != "https://cdn.example/audio-1.mp3" || second != "https://cdn.example/audio-1.mp3" {
		t.Fatalf("unexpected resolved URLs: %q and %q", first, second)
	}

	if calls.Load() != 1 {
		t.Fatalf("expected 1 upstream lookup, got %d", calls.Load())
	}
}

func TestResolveRedirectURL_UsesDatabaseFallback(t *testing.T) {
	cache := &testRedirectCache{}
	database := &testDatabase{redirects: map[int]string{1: "https://cdn.example/from-db.mp3"}}
	var calls atomic.Int32

	p := &Podcast{
		Cache:    cache,
		Database: database,
		Client: &testClient{resolveFunc: func(rawURL string) (string, error) {
			calls.Add(1)
			return "https://cdn.example/from-client.mp3", nil
		}},
	}

	resolvedURL, err := p.ResolveRedirectURL(1)
	if err != nil {
		t.Fatal(err)
	}

	if resolvedURL != "https://cdn.example/from-db.mp3" {
		t.Fatalf("resolved URL = %q, want %q", resolvedURL, "https://cdn.example/from-db.mp3")
	}

	if calls.Load() != 0 {
		t.Fatalf("expected 0 upstream lookups, got %d", calls.Load())
	}

	if cached, ok := cache.GetRedirectURL("1"); !ok || cached != "https://cdn.example/from-db.mp3" {
		t.Fatalf("memcache value = (%q, %v), want (%q, true)", cached, ok, "https://cdn.example/from-db.mp3")
	}
}

func TestResolveRedirectURL_PersistsClientResultInDatabase(t *testing.T) {
	cache := &testRedirectCache{}
	database := &testDatabase{}

	p := &Podcast{
		Cache:    cache,
		Database: database,
		Client: &testClient{resolveFunc: func(rawURL string) (string, error) {
			return "https://cdn.example/from-client.mp3", nil
		}},
	}

	resolvedURL, err := p.ResolveRedirectURL(7)
	if err != nil {
		t.Fatal(err)
	}

	if resolvedURL != "https://cdn.example/from-client.mp3" {
		t.Fatalf("resolved URL = %q, want %q", resolvedURL, "https://cdn.example/from-client.mp3")
	}

	database.mu.Lock()
	storedURL := database.storedRedirects[7]
	storeCalls := database.storeCalls
	database.mu.Unlock()

	if storedURL != "https://cdn.example/from-client.mp3" {
		t.Fatalf("stored URL = %q, want %q", storedURL, "https://cdn.example/from-client.mp3")
	}

	if storeCalls != 1 {
		t.Fatalf("store calls = %d, want %d", storeCalls, 1)
	}
}

func TestResolveRedirectURL_GlobalLockSerializesLookups(t *testing.T) {
	cache := &testRedirectCache{}
	var inFlight atomic.Int32
	var maxInFlight atomic.Int32

	p := &Podcast{
		Cache:    cache,
		Database: &testDatabase{},
		Client: &testClient{resolveFunc: func(rawURL string) (string, error) {
			current := inFlight.Add(1)
			for {
				max := maxInFlight.Load()
				if current <= max || maxInFlight.CompareAndSwap(max, current) {
					break
				}
			}

			time.Sleep(25 * time.Millisecond)
			inFlight.Add(-1)
			return "https://cdn.example/" + url.PathEscape(rawURL), nil
		}},
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 5)

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if _, err := p.ResolveRedirectURL(id); err != nil {
				errCh <- fmt.Errorf("id %d: %w", id, err)
			}
		}(i)
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatal(err)
	}

	if maxInFlight.Load() > 1 {
		t.Fatalf("expected at most one concurrent upstream lookup, got %d", maxInFlight.Load())
	}
}
