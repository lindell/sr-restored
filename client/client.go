package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/lindell/sr-restored/domain"
	"github.com/pkg/errors"
)

const userAgent = "Play/31330 (Android/34; google/Pixel 8 Pro; Scale/2.75)"

var baseURL *url.URL

func init() {
	baseURL, _ = url.Parse("https://app-api.sr.se/")
}

// Client is a client for the Sveriges Radio API.
type Client struct {
	HTTPClient *http.Client
	Cache      FileInfoCache
}

func NewClient(cache FileInfoCache) *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
		Cache:      cache,
	}
}

func (c *Client) GetProgram(ctx context.Context, id int, feedTypes []domain.FeedType) (domain.Program, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	program, err := c.getProgram(ctx, id)
	if err != nil {
		return domain.Program{}, errors.WithMessage(err, "could not fetch program from api")
	}

	program.Episodes, program.Hash, err = c.getEpisodes(ctx, id, feedTypes)
	if err != nil {
		return domain.Program{}, errors.WithMessage(err, "could not fetch episodes from api")
	}

	return program, nil
}

func (c *Client) getEpisodes(ctx context.Context, id int, feedTypes []domain.FeedType) ([]domain.Episode, []byte, error) {
	u := baseURL.JoinPath("episodes")
	q := u.Query()
	q.Add("programId", fmt.Sprint(id))
	q.Add("skip", "0")
	q.Add("take", fmt.Sprint(500))
	u.RawQuery = q.Encode()

	resp, err := c.fetch(ctx, http.MethodGet, u.String())
	if err != nil {
		return nil, nil, err
	}

	var listing EpisodeListing
	if err := json.NewDecoder(resp.Body).Decode(&listing); err != nil {
		return nil, nil, err
	}

	episodes := make([]domain.Episode, 0, len(listing.Items))
	for _, episode := range listing.Items {
		if converted, err := c.convertEpisode(episode, feedTypes); err == nil {
			episodes = append(episodes, converted)
		} else {
			slog.Error("could not convert episode",
				"error", err.Error(),
			)
		}
	}
	return episodes, listing.Hash(), nil
}

func (c *Client) getProgram(ctx context.Context, id int) (domain.Program, error) {
	u := baseURL.JoinPath("programs", fmt.Sprint(id))

	resp, err := c.fetch(ctx, http.MethodGet, u.String())
	if err != nil {
		return domain.Program{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return domain.Program{}, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var programInfo ProgramInfo
	if err := json.NewDecoder(resp.Body).Decode(&programInfo); err != nil {
		return domain.Program{}, errors.WithMessage(err, "could not decode JSON")
	}

	return convertProgram(programInfo), nil
}

func (c *Client) fetch(ctx context.Context, method string, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req = req.WithContext(ctx)
	res, err := c.HTTPClient.Do(req)
	if res.StatusCode >= 400 {
		return res, statusCodeError{
			msg:  fmt.Sprintf("unexpected status code %d", res.StatusCode),
			code: res.StatusCode,
		}
	}

	return res, err
}
