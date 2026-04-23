package podcast

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lindell/go-stderrs/stderrs"
	"github.com/lindell/sr-restored/domain"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type Podcast struct {
	Client   Client
	Cache    Cache
	Database Database

	redirectLookupLock sync.Mutex

	RSSUrl          *url.URL
	PodcastsFileURL *url.URL

	Now func() time.Time
}

type Client interface {
	GetProgram(ctx context.Context, id int, feedTypes []domain.FeedType) (domain.Program, error)
	ResolveRedirectURL(rawURL string) (string, error)
}

type Cache interface {
	StoreRSS(key string, rawRSS []byte)
	GetRSS(key string) ([]byte, bool)
	StoreHash(key string, hash []byte)
	GetHash(key string) ([]byte, bool)
	StoreRedirectURL(key string, redirectURL string)
	GetRedirectURL(key string) (string, bool)
}

type Database interface {
	InsertEpisodes(ctx context.Context, episodes []domain.Episode) error

	GetProgram(ctx context.Context, programID int) (domain.Program, error)
	InsertProgram(ctx context.Context, program domain.Program) error
	GetRedirectURL(ctx context.Context, audioFileID int) (string, error)
	StoreRedirectURL(ctx context.Context, audioFileID int, redirectURL string) error
}

func (p *Podcast) GetPodcast(ctx context.Context, id int, feedTypes []domain.FeedType) (rawRSS []byte, hash []byte, err error) {
	before := time.Now()
	cached := false

	cacheKey := fmt.Sprintf("%d:%v", id, feedTypes)

	defer func() {
		rssGetSecondsMetric.With(prometheus.Labels{
			"cached": fmt.Sprint(cached),
		}).Observe(time.Since(before).Seconds())

		rssGetTotalMetric.With(prometheus.Labels{
			"program_id": fmt.Sprint(id),
		}).Inc()
	}()

	if rss, ok := p.Cache.GetRSS(cacheKey); ok {
		cached = true

		hash, _ := p.Cache.GetHash(cacheKey)
		return rss, hash, nil
	}

	program, err := p.Client.GetProgram(ctx, id, feedTypes)
	if err != nil {
		// Try to fetch from DB as a backup
		var dbErr error
		program, dbErr = p.Database.GetProgram(ctx, id)
		if dbErr != nil {
			return nil, nil, errors.WithMessage(err, dbErr.Error())
		}
	} else {
		for i := range program.Episodes {
			program.Episodes[i].FileURL = p.toAudioProxyURL(program.Episodes[i].FileURL)
		}
	}

	rss := p.convertToPodRSS(program)

	raw, err := xml.Marshal(rss)
	if err != nil {
		return nil, nil, err
	}

	gzipedBuffer := bytes.NewBuffer(nil)
	gzipWriter := gzip.NewWriter(gzipedBuffer)
	if _, err := gzipWriter.Write(raw); err != nil {
		gzipWriter.Close()
		return nil, nil, errors.WithMessage(err, "could not gzip data")
	}
	gzipWriter.Close()
	gzipedRaw := gzipedBuffer.Bytes()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		p.Cache.StoreHash(cacheKey, program.Hash)
		p.Cache.StoreRSS(cacheKey, gzipedRaw)

		err := p.Database.InsertEpisodes(ctx, program.Episodes)
		if err != nil {
			slog.Error(err.Error())
		}

		err = p.Database.InsertProgram(ctx, program)
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	return gzipedRaw, program.Hash, nil
}

func (p *Podcast) IsNewestVersion(ctx context.Context, id int, feedTypes []domain.FeedType, hash []byte) (isNewest bool) {
	defer func() {
		hashLookup.With(prometheus.Labels{
			"success": fmt.Sprint(isNewest),
		}).Inc()
	}()

	cachedHash, ok := p.Cache.GetHash(fmt.Sprintf("%d:%v", id, feedTypes))
	if !ok {
		return false
	}

	return bytes.Equal(cachedHash, hash)
}

func (p *Podcast) ResolveRedirectURL(id int) (string, error) {
	cacheKey := fmt.Sprint(id)

	if cached, ok := p.Cache.GetRedirectURL(cacheKey); ok {
		return cached, nil
	}

	if cached, ok := p.getRedirectURLFromDatabase(id); ok {
		p.Cache.StoreRedirectURL(cacheKey, cached)
		return cached, nil
	}

	p.redirectLookupLock.Lock()
	defer p.redirectLookupLock.Unlock()

	resolvedURL, err := p.Client.ResolveRedirectURL(fmt.Sprintf("https://www.sverigesradio.se/topsy/ljudfil/%d", id))
	if err != nil {
		return "", errors.WithMessage(err, "could not resolve redirect URL")
	}

	p.Cache.StoreRedirectURL(cacheKey, resolvedURL)
	p.storeRedirectURLInDatabase(id, resolvedURL)

	return resolvedURL, nil
}

func (p *Podcast) getRedirectURLFromDatabase(id int) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	redirectURL, err := p.Database.GetRedirectURL(ctx, id)
	if err == nil {
		return redirectURL, true
	}

	if !stderrs.IsNotFound(err) {
		slog.Error("could not fetch redirect URL from db", "error", err)
	}

	return "", false
}

func (p *Podcast) storeRedirectURLInDatabase(id int, redirectURL string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := p.Database.StoreRedirectURL(ctx, id, redirectURL); err != nil {
		slog.Error("could not store redirect URL in db", "error", err)
	}
}

func (p *Podcast) toAudioProxyURL(rawURL string) string {
	if p.RSSUrl == nil {
		return rawURL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	if u.Scheme != "https" || u.Host != "www.sverigesradio.se" {
		return rawURL
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "topsy" || parts[1] != "ljudfil" {
		return rawURL
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		return rawURL
	}

	base := *p.RSSUrl
	base.Path = "/"
	base.RawPath = ""

	return base.JoinPath("audio-file", fmt.Sprint(id)).String()
}
