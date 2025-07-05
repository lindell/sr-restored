package httpserver

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/lindell/sr-restored/parallel"
	"github.com/lindell/sr-restored/podcast"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	Podcast         *podcast.Podcast
	parallelFetcher *parallel.SharedGetter[rssWithHash]
}

type rssWithHash struct {
	value []byte
	hash  []byte
}

// NewServer creates a new HTTP server with the given podcast service.
func NewServer(podcast *podcast.Podcast) *Server {
	return &Server{
		Podcast: podcast,
		parallelFetcher: parallel.NewSharedGetter(func(ctx context.Context, id int) (rssWithHash, error) {
			rawRSS, hash, err := podcast.GetPodcast(ctx, id)
			return rssWithHash{value: rawRSS, hash: hash}, err
		}),
	}
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

	rootMux.HandleFunc("/debug/pprof/", pprof.Index)
	rootMux.HandleFunc("/debug/pprof/{action}", pprof.Index)
	rootMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

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

	result, err := s.parallelFetcher.Fetch(r.Context(), int(id))
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Add("Content-Type", "application/xml")
	if len(result.hash) > 0 {
		if bytes.Equal(result.hash, sentHash) {
			responseWithCacheHit(w, result.hash)
			return
		}
		w.Header().Add("Etag", fmt.Sprintf("\"%s\"", base64.StdEncoding.EncodeToString(result.hash)))
	}
	respondWithGzippedData(w, r, result.value)
}

// The data is usually gzipped, but to support clients that do not support gzip,
// we check the Accept-Encoding header and decode the gzip data if necessary.
func respondWithGzippedData(w http.ResponseWriter, r *http.Request, data []byte) {
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
		return
	}

	slog.Debug("responding with non-gzipped data")
	gzipReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		slog.Error("failed to create gzip reader", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer gzipReader.Close()

	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, gzipReader)
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
