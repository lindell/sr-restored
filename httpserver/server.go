package httpserver

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/lindell/sr-restored/podcast"
)

type Server struct {
	Podcast *podcast.Podcast
}

func (s *Server) ListenAndServe(addr string) error {
	handler := s.Handler()

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		Handler:      handler,
		Addr:         addr,
	}

	return server.ListenAndServe()
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/rss/{id}", s.getRSS)

	mux.Handle("/", gziphandler.GzipHandler(http.FileServer(http.Dir("./static"))))

	return loggingMiddleware(slog.Default())(mux)
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
