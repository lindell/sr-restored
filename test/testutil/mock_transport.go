package testutil

import (
	"hash/fnv"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strconv"
	"strings"
)

type MockTransport struct {
	fileMocks map[string]string
}

func (mt *MockTransport) AddFileRespons(path string, filename string) {
	if mt.fileMocks == nil {
		mt.fileMocks = map[string]string{}
	}

	mt.fileMocks[path] = filename
}

func (mt *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	recorder := httptest.NewRecorder()

	if filename, ok := mt.fileMocks[req.URL.Path]; ok {
		bb, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		recorder.Write(bb)
		recorder.WriteHeader(http.StatusOK)
		return recorder.Result(), nil
	}

	if strings.HasSuffix(req.URL.Path, ".m4a") || strings.HasSuffix(req.URL.Path, ".mp4") {
		recorder.Header().Set("Content-Length", strconv.Itoa(hashFileSize(path.Base(req.URL.Path))))
		recorder.Header().Set("Content-Type", "audio/mp4")
		recorder.WriteHeader(http.StatusOK)
		return recorder.Result(), nil
	}

	recorder.WriteString("not found in mocks")
	recorder.WriteHeader(http.StatusNotFound)

	return recorder.Result(), nil
}

func hashFileSize(name string) int {
	h := fnv.New32a()
	h.Write([]byte(name))
	return int(h.Sum32())
}
