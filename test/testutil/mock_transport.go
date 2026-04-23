package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

type MockTransport struct {
	fileMocks map[string]string
	headMocks map[string]headResponse
}

type headResponse struct {
	contentType   string
	contentLength int
}

func (mt *MockTransport) AddFileRespons(path string, filename string) {
	if mt.fileMocks == nil {
		mt.fileMocks = map[string]string{}
	}

	mt.fileMocks[path] = filename
}

func (mt *MockTransport) AddHeadResponse(url string, contentType string, contentLength int) {
	if mt.headMocks == nil {
		mt.headMocks = map[string]headResponse{}
	}

	mt.headMocks[url] = headResponse{contentType: contentType, contentLength: contentLength}
}

func (mt *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	recorder := httptest.NewRecorder()

	if filename, ok := mt.fileMocks[req.URL.Path]; ok {
		bb, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		recorder.Write(bb)
	} else if req.Method == http.MethodHead {
		fullURL := req.URL.String()
		if hr, ok := mt.headMocks[fullURL]; ok {
			recorder.Header().Set("Content-Type", hr.contentType)
			recorder.Header().Set("Content-Length", fmt.Sprint(hr.contentLength))
		} else if strings.HasSuffix(req.URL.Path, ".mp3") {
			recorder.Header().Set("Content-Type", "audio/mpeg")
			recorder.Header().Set("Content-Length", "90000000")
		} else if strings.HasSuffix(req.URL.Path, ".m4a") {
			recorder.Header().Set("Content-Type", "audio/mp4")
			recorder.Header().Set("Content-Length", "300000000")
		} else {
			recorder.WriteString("not found in mocks")
			recorder.WriteHeader(http.StatusNotFound)
		}
	} else {
		recorder.WriteString("not found in mocks")
		recorder.WriteHeader(http.StatusNotFound)
	}

	return recorder.Result(), nil
}
