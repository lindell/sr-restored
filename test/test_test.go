package test_test

import (
	"context"
	"encoding/xml"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/lindell/sr-restored/podcast"
)

// TestLive requests a real RSS feed. This test might fail if Sveriges Radios API is down
func TestLive(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() { cancel() }()

	u := setup(ctx, t)

	rssURL := u.JoinPath("rss/2519").String()

	// Make first, not cached request
	resp, err := http.Get(rssURL)
	if err != nil {
		t.Fatal(err)
	}
	var rss podcast.RSS
	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		t.Fatal(err)
	}
	expectedTitle := "P3 Dokumentär — SR-restored"
	gotTitle := rss.Channel.Title
	if gotTitle != expectedTitle {
		t.Errorf("RSS title = %v, want %v", gotTitle, expectedTitle)
	}

	// Make second, cached request
	beforeReq := time.Now()
	respCached, err := http.Get(rssURL)
	if err != nil {
		t.Fatal(err)
	}
	requestDuration := time.Since(beforeReq)
	if requestDuration > 5*time.Millisecond {
		t.Fatalf("cached request took too long: %s", requestDuration)
	}
	var cachedRSS podcast.RSS
	if err := xml.NewDecoder(respCached.Body).Decode(&cachedRSS); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(cachedRSS, rss) {
		t.Fatal("cached RSS does not match first returned one")
	}

	_ = respCached
}
