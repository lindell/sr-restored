package podcast

import (
	"fmt"
	"strings"
	"time"

	"github.com/lindell/sr-restored/domain"
)

func (p *Podcast) convertToPodRSS(program domain.Program) RSS {
	rss := baseRSS

	title := fmt.Sprintf("%s — SR-restored", program.Name)
	selfURL := p.RSSUrl.JoinPath(fmt.Sprint(program.ID))

	rss.Channel.Title = title

	rss.Channel.ItunesNewFeedURL = selfURL.String()
	rss.Channel.AtomLink.Href = selfURL.String()
	rss.Channel.AtomLink.Rel = "self"
	rss.Channel.AtomLink.Type = "application/rss+xml"
	rss.Channel.LastBuildDate = p.Now().Format(time.RFC1123)

	rss.Channel.Image.URL = program.ImageURL
	rss.Channel.Image.Title = title
	rss.Channel.Image.Link = program.URL

	rss.Channel.ItunesImage = program.ImageURL
	rss.Channel.ItunesSummary = program.Description
	rss.Channel.ItunesAuthor = "Sveriges Radio"
	// rss.Channel.ItunesCategory

	rss.Channel.ItunesOwner.Name = title
	rss.Channel.ItunesOwner.Email = "johan+srrestored@lindell.me"

	rss.Channel.Link = program.URL
	rss.Channel.Description = program.Description
	rss.Channel.Language = "sv"
	rss.Channel.Copyright = program.Copyright

	for _, episode := range program.Episodes {
		rss.Channel.Item = append(rss.Channel.Item, convertEpisode(episode))
	}

	return rss
}

func convertEpisode(original domain.Episode) PodItem {
	var target PodItem

	description := fmt.Sprintf(`<p>%s</p>
		<hr>
		<p>
			Innehållet i denna podcast kommer från Sveriges Radio. Podcastflödet har tillgängliggjorts av <a href="https://sr-restored.se">SR-restored</a> då Sveriges Radio ej publicerar alla poddar i sina RSS-flöden. Läs mer på <a href="https://sr-restored.se">sr-restored.se</a>.
		</p>
	`, original.Description)

	target.Title = original.Title
	target.Description = description

	// URL and GUID persistence
	target.Link = original.URL
	target.Guid.Text = fmt.Sprintf("rss:sr.se/pod/eid/%d", original.ID)
	target.Guid.IsPermaLink = "false"

	target.PubDate = original.PublishDate.Format(time.RFC1123)

	target.Programid = original.ID

	target.Summary = description

	target.Author = "Sveriges Radio"

	target.Keywords = strings.ReplaceAll(original.Title, " ", ",")

	target.Image.Href = original.ImageURL

	target.Duration = fmtDuration(
		time.Second * time.Duration(original.FileDurationSeconds),
	)
	target.Subtitle = original.Description

	target.Enclosure.URL = original.FileURL
	target.Enclosure.Length = original.FileBytes
	target.Enclosure.Type = original.ContentType

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
