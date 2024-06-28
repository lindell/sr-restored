package testutil

import (
	"net/http"
	"net/http/httptest"
	"os"
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
	} else {
		recorder.WriteString("not found in mocks")
		recorder.WriteHeader(http.StatusNotFound)
	}

	return recorder.Result(), nil
}
