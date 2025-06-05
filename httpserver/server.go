package httpserver

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/lindell/sr-restored/podcast"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	Podcast *podcast.Podcast
}

func (s *Server) ListenAndServe(ctx context.Context, addr string) (stop func(context.Context) error) {
	handler := s.Handler()

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      handler,
		Addr:         addr,
	}

	go func() {
		slog.Info(fmt.Sprintf("listening on %s", addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server closed",
				"error", err,
			)
		}
	}()

	return func(ctx context.Context) error {
		slog.Info("shutting down http server")
		return server.Shutdown(ctx)
	}
}

func (s *Server) Handler() http.Handler {
	rootMux := http.NewServeMux()

	rootMux.Handle("/metrics", promhttp.Handler())

	mainMux := http.NewServeMux()
	rootMux.Handle("/", loggingMiddleware(slog.Default())(gziphandler.GzipHandler(mainMux)))

	mainMux.HandleFunc("/rss/{id}", s.getRSS)
	mainMux.Handle("/", http.FileServer(http.Dir("./static")))

	return rootMux
}

func (s *Server) getRSS(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	sentHash := getIfNoneMatchHash(r)
	if sentHash != nil && s.Podcast.IsNewestVersion(r.Context(), int(id), sentHash) {
		responseWithCacheHit(w, sentHash)
		return
	}

	b, hash, err := s.Podcast.GetPodcast(r.Context(), int(id))
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Add("Content-Type", "application/xml")
	if len(hash) > 0 {
		if bytes.Equal(hash, sentHash) {
			responseWithCacheHit(w, hash)
			return
		}
		w.Header().Add("Etag", fmt.Sprintf("\"%s\"", base64.StdEncoding.EncodeToString(hash)))
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func responseWithCacheHit(w http.ResponseWriter, hash []byte) {
	w.Header().Add("Etag", fmt.Sprintf("\"%s\"", base64.StdEncoding.EncodeToString(hash)))
	w.WriteHeader(http.StatusNotModified)
}

func getIfNoneMatchHash(r *http.Request) []byte {
	etag := r.Header.Get("If-None-Match")
	if etag == "" {
		return nil
	}

	b64hash, err := strconv.Unquote(etag)
	if err != nil {
		return nil
	}

	hash, err := base64.StdEncoding.DecodeString(b64hash)
	if err != nil {
		return nil
	}

	return hash
}
