package httpserver

import (
	"context"
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
	rootMux.Handle("/", loggingMiddleware(slog.Default())(mainMux))

	mainMux.HandleFunc("/rss/{id}", s.getRSS)
	mainMux.Handle("/", gziphandler.GzipHandler(http.FileServer(http.Dir("./static"))))

	return rootMux
}

func (s *Server) getRSS(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	b, err := s.Podcast.GetPodcast(r.Context(), int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}
