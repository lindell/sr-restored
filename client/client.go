package client

import (
	"context"
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/lindell/sr-restored/domain"
	"github.com/pkg/errors"
)

var baseURL *url.URL

func init() {
	baseURL, _ = url.Parse("https://api.sr.se/api/v2/")
}

func GetProgram(ctx context.Context, id int) (domain.Program, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	program, err := getProgram(ctx, id)
	if err != nil {
		return domain.Program{}, errors.WithMessage(err, "could not fetch program from api")
	}

	program.Episodes, err = getEpisodes(ctx, id)
	if err != nil {
		return domain.Program{}, errors.WithMessage(err, "could not fetch episodes from api")
	}

	return program, nil
}

func getEpisodes(ctx context.Context, id int) ([]domain.Episode, error) {
	u := baseURL.JoinPath("episodes/index?audioquality=hi&size=500")
	q := u.Query()
	q.Add("programid", fmt.Sprint(id))
	u.RawQuery = q.Encode()

	resp, err := fetch(ctx, http.MethodGet, u.String())
	if err != nil {
		return nil, err
	}

	var listing EpisodeListing
	if err := xml.NewDecoder(resp.Body).Decode(&listing); err != nil {
		return nil, err
	}

	episodes := make([]domain.Episode, 0, len(listing.Episodes.Episode))
	for _, episode := range listing.Episodes.Episode {
		if converted, err := convertEpisode(episode); err == nil {
			episodes = append(episodes, converted)
		} else {
			slog.Error("could not convert episode",
				"error", err.Error(),
			)
		}
	}
	return episodes, nil
}

func getProgram(ctx context.Context, id int) (domain.Program, error) {
	u := baseURL.JoinPath("programs", fmt.Sprint(id))

	resp, err := fetch(ctx, http.MethodGet, u.String())
	if err != nil {
		return domain.Program{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return domain.Program{}, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var programInfo ProgramInfo
	if err := xml.NewDecoder(resp.Body).Decode(&programInfo); err != nil {
		return domain.Program{}, errors.WithMessage(err, "could not decode XML")
	}

	return convertProgram(programInfo), nil
}

func fetch(ctx context.Context, method string, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode >= 400 {
		return res, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	return res, err
}
