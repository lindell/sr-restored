package main

import (
	"context"
	"os"

	"github.com/lindell/sr-restored/run"
)

func main() {
	err := run.Run(context.Background(), run.Config{
		ServerAddr:  ":8080",
		BaseURL:     getEnv("BASE_URL", "http://localhost:8080"),
		PostgresURL: getEnv("DATABASE_URL", ""),
	})
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
