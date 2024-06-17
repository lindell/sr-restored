package main

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/lindell/sr-restored/httpserver"
	"github.com/lindell/sr-restored/memcache"
	"github.com/lindell/sr-restored/podcast"
	"github.com/lindell/sr-restored/postgres"
)

const addr = ":8080"

func main() {
	// The URL in which the service is hosted
	baseURLStr := getEnv("BASE_URL", "http://localhost:8080")
	postgresURL := getEnv("DATABASE_URL", "")

	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		panic(err)
	}

	postgresDB, err := postgres.NewDB(postgres.Config{
		PostgresURL: postgresURL,
	})
	if err != nil {
		panic(err)
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
		panic(err)
	}
}

func getEnv(name string, defaultValue string) string {
	val, hasVal := os.LookupEnv(name)
	if !hasVal {
		return defaultValue
	}
	return val
}
