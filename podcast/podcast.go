package podcast

import (
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
}

type Database interface {
	InsertEpisodes(ctx context.Context, episodes []domain.Episode) error

	GetProgram(ctx context.Context, programID int) (domain.Program, error)
	InsertProgram(ctx context.Context, program domain.Program) error
}

func (p *Podcast) GetPodcast(ctx context.Context, id int) ([]byte, error) {
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
		return rss, nil
	}

	program, err := client.GetProgram(ctx, id)
	if err != nil {
		// Try to fetch from DB as a backup
		var dbErr error
		program, dbErr = p.Database.GetProgram(ctx, id)
		if dbErr != nil {
			return nil, errors.WithMessage(err, dbErr.Error())
		}
	}

	rss := p.convertToPodRSS(program)

	raw, err := xml.MarshalIndent(rss, "  ", "    ")
	if err != nil {
		return nil, err
	}

	go func() {
		p.Cache.StoreRSS(id, raw)

		err := p.Database.InsertEpisodes(context.Background(), program.Episodes)
		if err != nil {
			slog.Error(err.Error())
		}

		err = p.Database.InsertProgram(context.Background(), program)
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	return raw, nil
}
