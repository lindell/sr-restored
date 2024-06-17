package podcast

import (
	"context"
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/lindell/sr-restored/client"
	"github.com/lindell/sr-restored/domain"
	"github.com/pkg/errors"
)

type Podcast struct {
	Cache    Cache
	Database Database

	RSSUrl *url.URL
}

type Cache interface {
	StoreRSS(id int, rawRSS []byte)
	GetRSS(id int) ([]byte, bool)
}

type Database interface {
	InsertEpisodes(ctx context.Context, episodes []domain.Episode) error

	GetProgram(ctx context.Context, programID int) (domain.Program, error)
	InsertProgram(ctx context.Context, program domain.Program) error
}

func (p *Podcast) GetPodcast(ctx context.Context, id int) ([]byte, error) {
	if rss, ok := p.Cache.GetRSS(id); ok {
		return rss, nil
	}

	program, err := client.GetProgram(ctx, id)
	if err != nil {
		// Try to fetch from DB as a backup
		var dbErr error
		program, dbErr = p.Database.GetProgram(ctx, id)
		if dbErr != nil {
			return nil, errors.WithMessage(err, dbErr.Error())
		}
	}

	rss := baseRSS

	title := fmt.Sprintf("%s — SR-restored", program.Name)
	selfURL := p.RSSUrl.JoinPath(fmt.Sprint(program.ID))

	rss.Channel.Title = title

	rss.Channel.ItunesNewFeedURL = selfURL.String()
	rss.Channel.AtomLink.Href = selfURL.String()
	rss.Channel.AtomLink.Rel = "self"
	rss.Channel.AtomLink.Type = "application/rss+xml"
	rss.Channel.LastBuildDate = time.Now().Format(time.RFC1123) // TODO

	rss.Channel.Image.URL = program.ImageURL
	rss.Channel.Image.Title = title
	rss.Channel.Image.Link = program.URL

	rss.Channel.ItunesImage = program.ImageURL
	rss.Channel.ItunesSummary = program.Description
	rss.Channel.ItunesAuthor = "Sveriges Radio"
	// rss.Channel.ItunesCategory

	rss.Channel.ItunesOwner.Name = title
	rss.Channel.ItunesOwner.Email = program.Email

	rss.Channel.Link = program.URL
	rss.Channel.Description = program.Description
	rss.Channel.Language = "sv"
	rss.Channel.Copyright = program.Copyright

	for _, episode := range program.Episodes {
		rss.Channel.Item = append(rss.Channel.Item, convertEpisode(episode))
	}

	raw, err := xml.Marshal(rss)
	if err != nil {
		return nil, err
	}

	go func() {
		p.Cache.StoreRSS(id, raw)

		err := p.Database.InsertEpisodes(context.Background(), program.Episodes)
		if err != nil {
			slog.Error(err.Error())
		}

		err = p.Database.InsertProgram(context.Background(), program)
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	return raw, nil
}

func convertEpisode(original domain.Episode) PodItem {
	var target PodItem
	target.Title = original.Title
	target.Description = original.Description

	// URL and GUID persistence
	target.Link = original.URL
	target.Guid.Text = fmt.Sprintf("rss:sr.se/pod/eid/%d", original.ID)
	target.Guid.IsPermaLink = "false"

	target.PubDate = original.PublishDate.Format(time.RFC1123)

	target.Programid = original.ID

	target.Summary = original.Description

	target.Author = "Sveriges Radio"

	target.Keywords = strings.ReplaceAll(original.Title, " ", ",")

	target.Image.Href = original.ImageURL

	target.Duration = fmtDuration(
		time.Second * time.Duration(original.FileDurationSeconds),
	)
	target.Subtitle = original.Description

	target.Enclosure.URL = original.FileURL
	target.Enclosure.Length = original.FileBytes
	target.Enclosure.Type = "audio/mpeg" // TODO: Determine based on

	return target
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)

	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
