package testutil

import (
	"net/http"
	"net/http/httptest"
	"os"
)

type MockTransport struct {
	fileMocks map[string]string
}

var mediaMocks = map[string]struct {
	length string
	cType  string
}{
	"/ljudit/p3/morgonpasset_i_p3/2026/02/morgonpasset_i_p3_20260218_0559_192.m4a": {"306219605", "audio/mp4"},
	"/autorec/et2w/p3/morgonpasset_i_p3/2026/02/SRP3_2026-02-17_055900_12660_a32.m4a": {"51829181", "audio/mp4"},
	"/ljudit/p3/morgonpasset_i_p3/2026/02/morgonpasset_i_p3_20260216_0559_192.m4a": {"306218152", "audio/mp4"},
	"/ljudit/p3/morgonpasset_i_p3/2026/02/morgonpasset_i_p3_20260215_0800_192.m4a": {"261229667", "audio/mp4"},
	"/ljudit/p3/morgonpasset_i_p3/2026/02/morgonpasset_i_p3_20260214_0800_192.m4a": {"261229511", "audio/mp4"},
	"/ljudit/p3/morgonpasset_i_p3/2026/02/morgonpasset_i_p3_20260213_0559_192.m4a": {"306218197", "audio/mp4"},
	"/ljudit/p3/morgonpasset_i_p3/2026/02/morgonpasset_i_p3_20260212_0559_192.m4a": {"306218431", "audio/mp4"},
	"/ljudit/p3/morgonpasset_i_p3/2026/02/morgonpasset_i_p3_20260211_0559_192.m4a": {"306219126", "audio/mp4"},
	"/ljudit/p3/morgonpasset_i_p3/2026/02/morgonpasset_i_p3_20260210_0559_192.m4a": {"306219241", "audio/mp4"},
	"/ljudit/p3/morgonpasset_i_p3/2026/02/morgonpasset_i_p3_20260209_0559_32.m4a": {"51828993", "audio/mp4"},
}

func (mt *MockTransport) AddFileRespons(path string, filename string) {
	if mt.fileMocks == nil {
		mt.fileMocks = map[string]string{}
	}

	mt.fileMocks[path] = filename
}

func (mt *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	recorder := httptest.NewRecorder()

	if info, ok := mediaMocks[req.URL.Path]; ok {
		recorder.Header().Set("Content-Length", info.length)
		recorder.Header().Set("Content-Type", info.cType)
		recorder.WriteHeader(http.StatusOK)
		return recorder.Result(), nil
	}

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
