package run

import (
	"fmt"
	"log/slog"
	"net/url"

	"github.com/lindell/sr-restored/httpserver"
	"github.com/lindell/sr-restored/memcache"
	"github.com/lindell/sr-restored/podcast"
	"github.com/lindell/sr-restored/postgres"
)

const addr = ":8080"

type Config struct {
	// The URL in which the service is hosted
	BaseURL     string
	PostgresURL string
}

func Run(config Config) error {
	baseURL, err := url.Parse(config.BaseURL)
	if err != nil {
		return err
	}

	postgresDB, err := postgres.NewDB(postgres.Config{
		PostgresURL: config.PostgresURL,
	})
	if err != nil {
		return err
	}

	cache := memcache.NewCache()

	podcast := &podcast.Podcast{
		Cache:    cache,
		Database: postgresDB,
		RSSUrl:   baseURL.JoinPath("rss"),
	}

	httpServer := httpserver.Server{
		Podcast: podcast,
	}

	slog.Info(fmt.Sprintf("listening on %s", addr))
	err = httpServer.ListenAndServe(addr)
	if err != nil {
		return err
	}

	return nil
}
