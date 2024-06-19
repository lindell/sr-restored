package scraper

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/lindell/sr-restored/domain"
	"github.com/pkg/errors"
)

type Scraper struct{}

var atLeastOneNewline = regexp.MustCompile(`[ \t\n]*\n[ \t\n]*`)

var baseURL *url.URL

func init() {
	baseURL, _ = url.Parse("https://sverigesradio.se/")
}

func (sc *Scraper) FindEpisodes(programID int) ([]domain.Episode, error) {
	doc, err := sc.fetchEpisodeHTML("https://sverigesradio.se/default.aspx?programid=2519")
	if err != nil {
		return nil, err
	}

	episodes, err := docMap(doc.Find(".episode-list-item"), func(s *goquery.Selection) (domain.Episode, error) {
		strID, exist := s.Attr("data-spa-item-id")
		if !exist {
			return domain.Episode{}, errors.New("could not find id")
		}
		id, err := strconv.Atoi(strID)
		if err != nil {
			return domain.Episode{}, errors.WithMessage(err, "could not parse episode id")
		}

		headerSel := s.Find(".audio-heading__title .heading")
		if headerSel.Size() != 1 {
			return domain.Episode{}, fmt.Errorf("expected one header, found %d", headerSel.Size())
		}
		title := headerSel.Text()

		timeSel := s.Find("time")
		datetimeStr, ok := timeSel.Attr("datetime")
		if !ok {
			return domain.Episode{}, errors.New("could not find datetime")
		}
		publishDate, err := time.Parse("2006-01-02 15:04:05Z", datetimeStr)
		if err != nil {
			return domain.Episode{}, errors.WithMessage(err, "could not parse datetime")
		}

		description := s.Find(".episode-list-item__description.is-hidden-touch").Text()
		if description == "" {
			return domain.Episode{}, errors.New("could not find description")
		}
		description = strings.TrimSpace(description)
		description = atLeastOneNewline.ReplaceAllString(description, "\n")

		episodeURL, ok := s.Find(".audio-heading a").Attr("href")
		if !ok {
			return domain.Episode{}, errors.New("could not find url")
		}
		episodeURL = baseURL.JoinPath(episodeURL).String()

		fileURL, ok := s.Find("[data-stat-metadata-id]").Attr("href")
		if !ok {
			return domain.Episode{}, errors.New("could not find file url")
		}
		if strings.HasPrefix(fileURL, "//") {
			fileURL = "https:" + fileURL
		}

		durationStr, ok := s.Find(".audio-heading__meta abbr").Attr("title")
		if !ok {
			return domain.Episode{}, errors.New("could not find duration")
		}
		fileDuration, err := parseDuring(durationStr)
		if err != nil {
			return domain.Episode{}, err
		}

		imageURL, ok := s.Find(".image img").Attr("data-src")
		if !ok {
			return domain.Episode{}, errors.New("could not find image src url")
		}
		imageU, err := url.Parse(imageURL)
		if err != nil {
			return domain.Episode{}, errors.WithMessage(err, "could not parse image src url")
		}
		q := imageU.Query()
		q.Del("preset")
		imageU.RawQuery = q.Encode()

		contentType, fileSize, err := sc.getFileInfo(fileURL)
		if err != nil {
			return domain.Episode{}, err
		}

		return domain.Episode{
			ID:          id,
			ProgramID:   programID,
			Title:       title,
			Description: description,
			URL:         episodeURL,
			PublishDate: publishDate,
			ImageURL:    imageU.String(),

			ContentType:         contentType,
			FileURL:             fileURL,
			FileDurationSeconds: int(fileDuration / time.Second),
			FileBytes:           fileSize,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return episodes, nil
}

func (s *Scraper) fetchEpisodeHTML(url string) (*goquery.Document, error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *Scraper) getFileInfo(fileURL string) (contentType string, size int, err error) {
	res, err := http.Head(fileURL)
	if err != nil {
		return "", 0, errors.WithMessage(err, "could not fetch file url")
	}

	contentLength := res.Header.Get("Content-Length")
	size, err = strconv.Atoi(contentLength)
	if err != nil {
		return "", 0, errors.WithMessage(err, "could not parse file url content length")
	}

	return res.Header.Get("Content-Type"), size, nil
}
