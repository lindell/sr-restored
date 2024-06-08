package podcast

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/lindell/sr-uncensored/client"
	"github.com/pkg/errors"
)

type Podcast struct {
	Cache Cache

	RSSUrl *url.URL
}

type Cache interface {
	StoreRSS(id int, rawRSS []byte)
	GetRSS(id int) ([]byte, bool)
}

func (p *Podcast) GetPodcast(ctx context.Context, id int) ([]byte, error) {
	if rss, ok := p.Cache.GetRSS(id); ok {
		return rss, nil
	}

	program, err := client.GetProgram(ctx, id)
	if err != nil {
		return nil, errors.WithMessage(err, "could not get program data")
	}

	episodes, err := client.GetEpisodes(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = episodes

	rss := baseRSS

	title := fmt.Sprintf("%s (sr-uncensored)", program.Program.Name)
	selfURL := p.RSSUrl.JoinPath(program.Program.ID)

	rss.Channel.Title = title

	rss.Channel.ItunesNewFeedURL = selfURL.String()
	rss.Channel.AtomLink.Href = selfURL.String()
	rss.Channel.AtomLink.Rel = "self"
	rss.Channel.AtomLink.Type = "application/rss+xml"
	rss.Channel.LastBuildDate = time.Now().Format(time.RFC1123) // TODO

	rss.Channel.Image.URL = program.Program.Programimage
	rss.Channel.Image.Title = title
	rss.Channel.Image.Link = program.Program.Programurl

	rss.Channel.ItunesImage = program.Program.Programimage
	rss.Channel.ItunesSummary = program.Program.Description
	rss.Channel.ItunesAuthor = "Sveriges Radio"
	// rss.Channel.ItunesCategory

	rss.Channel.ItunesOwner.Name = title
	rss.Channel.ItunesOwner.Email = program.Program.Email

	rss.Channel.Link = program.Program.Programurl
	rss.Channel.Description = program.Program.Description
	rss.Channel.Language = "sv"
	rss.Channel.Copyright = program.Copyright

	for _, episode := range episodes.Episodes.Episode {
		rss.Channel.Item = append(rss.Channel.Item, convertEpisode(episode))
	}

	raw, err := xml.Marshal(rss)
	if err != nil {
		return nil, err
	}

	p.Cache.StoreRSS(id, raw)

	return raw, nil
}

func convertEpisode(original client.Episode) PodItem {
	var target PodItem
	target.Title = original.Title
	target.Description = original.Description

	// URL and GUID persistence
	target.Link = original.URL
	target.Guid.Text = original.URL
	target.Guid.IsPermaLink = "true"

	target.PubDate = original.Publishdateutc

	target.Programid = original.Program.ID
	// target.Poddid

	target.Summary = original.Description

	target.Author = "Sveriges Radio"

	target.Keywords = strings.ReplaceAll(original.Title, " ", ",")

	target.Image.Href = original.Imageurl

	target.Duration = original.Downloadpodfile.Duration
	target.Subtitle = original.Description

	target.Enclosure.URL = original.Downloadpodfile.URL
	target.Enclosure.Length = original.Downloadpodfile.Filesizeinbytes
	target.Enclosure.Type = "audio/mpeg" // TODO: Determine based on

	return target
}
