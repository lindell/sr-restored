package httpserver

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Info(
						"err", err,
						"trace", debug.Stack(),
					)
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			logger.Info("http request",
				"status", wrapped.status,
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"user-agent", r.Header.Get("User-Agent"),
				"headers", flattenHeaders(r.Header),
				"duration", time.Since(start),
			)
		}

		return http.HandlerFunc(fn)
	}
}

// flattenHeaders takes a headers map and converts it to only one value per header name
func flattenHeaders(headers map[string][]string) map[string]string {
	logHeaders := map[string]string{}

	for name, values := range headers {
		logHeaders[name] = strings.Join(values, " ||| ")
	}

	return logHeaders
}
