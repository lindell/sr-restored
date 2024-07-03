package run

import (
	"context"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lindell/sr-restored/httpserver"
	"github.com/lindell/sr-restored/memcache"
	"github.com/lindell/sr-restored/podcast"
	"github.com/lindell/sr-restored/postgres"
)

type Config struct {
	// Which address to listen on (HTTP)
	ServerAddr string

	// The URL in which the service is hosted
	BaseURL     string
	PostgresURL string

	Now func() time.Time
}

func Run(ctx context.Context, config Config) error {
	if config.Now == nil {
		config.Now = time.Now
	}

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

		Now: config.Now,
	}

	httpServer := httpserver.Server{
		Podcast: podcast,
	}

	stopHTTPServer := httpServer.ListenAndServe(ctx, config.ServerAddr)

	waitForStop(ctx)

	stopHTTPServer(ctx)

	return nil
}

func waitForStop(ctx context.Context) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	select {
	case <-sigc:
	case <-ctx.Done():
	}
}
