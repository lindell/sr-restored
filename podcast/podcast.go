package podcast

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/lindell/sr-restored/client"
	"github.com/lindell/sr-restored/domain"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type Podcast struct {
	Cache    Cache
	Database Database

	RSSUrl *url.URL

	Now func() time.Time
}

type Cache interface {
	StoreRSS(id int, rawRSS []byte)
	GetRSS(id int) ([]byte, bool)
	StoreHash(id int, hash []byte)
	GetHash(id int) ([]byte, bool)
}

type Database interface {
	InsertEpisodes(ctx context.Context, episodes []domain.Episode) error

	GetProgram(ctx context.Context, programID int) (domain.Program, error)
	InsertProgram(ctx context.Context, program domain.Program) error
}

func (p *Podcast) GetPodcast(ctx context.Context, id int) (rawRSS []byte, hash []byte, err error) {
	before := time.Now()
	cached := false
	defer func() {
		rssGetSecondsMetric.With(prometheus.Labels{
			"cached": fmt.Sprint(cached),
		}).Observe(time.Since(before).Seconds())

		rssGetTotalMetric.With(prometheus.Labels{
			"program_id": fmt.Sprint(id),
		}).Inc()
	}()

	if rss, ok := p.Cache.GetRSS(id); ok {
		cached = true

		hash, _ := p.Cache.GetHash(id)
		return rss, hash, nil
	}

	program, err := client.GetProgram(ctx, id)
	if err != nil {
		// Try to fetch from DB as a backup
		var dbErr error
		program, dbErr = p.Database.GetProgram(ctx, id)
		if dbErr != nil {
			return nil, nil, errors.WithMessage(err, dbErr.Error())
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
		return nil, nil, errors.WithMessage(err, "could not gzip data")
	}
	gzipWriter.Close()
	gzipedRaw := gzipedBuffer.Bytes()

	go func() {
		p.Cache.StoreHash(id, program.Hash)
		p.Cache.StoreRSS(id, gzipedRaw)

		err := p.Database.InsertEpisodes(context.Background(), program.Episodes)
		if err != nil {
			slog.Error(err.Error())
		}

		err = p.Database.InsertProgram(context.Background(), program)
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	return gzipedRaw, program.Hash, nil
}

func (p *Podcast) IsNewestVersion(ctx context.Context, id int, hash []byte) (isNewest bool) {
	defer func() {
		hashLookup.With(prometheus.Labels{
			"success": fmt.Sprint(isNewest),
		}).Inc()
	}()

	cachedHash, ok := p.Cache.GetHash(id)
	if !ok {
		return false
	}

	return bytes.Equal(cachedHash, hash)
}
