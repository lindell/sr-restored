package httpserver

import (
	"errors"
	"net/http"
)

func handleError(w http.ResponseWriter, _ *http.Request, err error) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(errorStatusCode(err))
	_, _ = w.Write([]byte(err.Error()))
}

func errorStatusCode(err error) int {
	var statusErr interface {
		HTTPStatusCode() int
	}

	if ok := errors.As(err, &statusErr); ok {
		return statusErr.HTTPStatusCode()
	}

	return http.StatusInternalServerError
}
