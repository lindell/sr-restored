package test_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/lindell/sr-restored/client"
	"github.com/lindell/sr-restored/test/testutil"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/golden"
)

func TestConvert(t *testing.T) {
	mockTransport := testutil.MockTransport{}
	prevClient := client.DefaultClient
	client.DefaultClient = &http.Client{
		Transport: &mockTransport,
	}
	defer func() { client.DefaultClient = prevClient }()
	mockTransport.AddFileRespons("/api/v2/programs/2519", "data/program.xml")
	mockTransport.AddFileRespons("/api/v2/episodes/index", "data/episodes.xml")

	ctx, cancel := context.WithCancel(context.Background())
	defer func() { cancel() }()

	u := setup(ctx, t)

	rssURL := u.JoinPath("rss/2519").String()

	resp, err := http.Get(rssURL)
	if err != nil {
		t.Fatal(err)
	}
	resultBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	golden.Assert(t, string(resultBody), "resulting-rss.xml")
}
