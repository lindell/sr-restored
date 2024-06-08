package client

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

var baseURL *url.URL

func init() {
	baseURL, _ = url.Parse("https://api.sr.se/api/v2/")
}

func GetEpisodes(ctx context.Context, id int) (EpisodeListing, error) {
	u := baseURL.JoinPath("episodes/index?audioquality=hi&size=1000")
	q := u.Query()
	q.Add("programid", fmt.Sprint(id))
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return EpisodeListing{}, err
	}

	var listing EpisodeListing
	if err := xml.NewDecoder(resp.Body).Decode(&listing); err != nil {
		return EpisodeListing{}, err
	}

	return listing, nil
}

func GetProgram(ctx context.Context, id int) (ProgramInfo, error) {
	u := baseURL.JoinPath("programs", fmt.Sprint(id))

	resp, err := http.Get(u.String())
	if err != nil {
		return ProgramInfo{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return ProgramInfo{}, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var programInfo ProgramInfo
	if err := xml.NewDecoder(resp.Body).Decode(&programInfo); err != nil {
		return ProgramInfo{}, errors.WithMessage(err, "could not decode XML")
	}

	return programInfo, nil
}
