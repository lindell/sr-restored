package podcast

import "encoding/xml"

type RSS struct {
	XMLName    xml.Name `xml:"rss"`
	Text       string   `xml:",chardata"`
	Itunes     string   `xml:"xmlns:itunes,attr"`
	Googleplay string   `xml:"xmlns:googleplay,attr"`
	Atom       string   `xml:"xmlns:atom,attr"`
	Sr         string   `xml:"xmlns:sr,attr"`
	Media      string   `xml:"xmlns:media,attr"`
	Version    string   `xml:"version,attr"`
	Channel    struct {
		Text             string `xml:",chardata"`
		ItunesNewFeedURL string `xml:"itunes:new-feed-url"`
		AtomLink         struct {
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"atom:link"`
		LastBuildDate string `xml:"lastBuildDate"`
		Image         struct {
			Text  string `xml:",chardata"`
			Title string `xml:"title"`
			Link  string `xml:"link"`
			URL   string `xml:"url"`
		} `xml:"image"`
		ItunesImage    string `xml:"itunes:image"`
		ItunesExplicit string `xml:"itunes:explicit"`
		ItunesSummary  string `xml:"itunes:summary"`
		ItunesAuthor   string `xml:"itunes:author"`
		ItunesCategory struct {
			Text     string `xml:",chardata"`
			AttrText string `xml:"text,attr"`
		} `xml:"category"`
		ItunesOwner struct {
			Text  string `xml:",chardata"`
			Name  string `xml:"name"`
			Email string `xml:"email"`
		} `xml:"owner"`
		Title       string    `xml:"title"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Copyright   string    `xml:"copyright"`
		Link        string    `xml:"link"`
		Item        []PodItem `xml:"item"`
	} `xml:"channel"`
}

type PodItem struct {
	Text        string `xml:",chardata"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Guid        struct {
		Text        string `xml:",chardata"`
		IsPermaLink string `xml:"isPermaLink,attr"`
	} `xml:"guid"`
	PubDate   string `xml:"pubDate"`
	Programid int    `xml:"sr:programid"`
	Poddid    string `xml:"sr:poddid"`
	Summary   string `xml:"itunes:summary"`
	Author    string `xml:"itunes:author"`
	Keywords  string `xml:"itunes:keywords"`
	Image     struct {
		Text string `xml:",chardata"`
		Href string `xml:"href,attr"`
	} `xml:"itunes:image"`
	Duration  string `xml:"itunes:duration"`
	Subtitle  string `xml:"itunes:subtitle"`
	Enclosure struct {
		Text   string `xml:",chardata"`
		URL    string `xml:"url,attr"`
		Length int    `xml:"length,attr"`
		Type   string `xml:"type,attr"`
	} `xml:"enclosure"`
}

var baseRSS RSS = RSS{
	Itunes:     "http://www.itunes.com/dtds/podcast-1.0.dtd",
	Googleplay: "http://www.google.com/schemas/play-podcasts/1.0",
	Atom:       "http://www.w3.org/2005/Atom",
	Sr:         "http://www.sverigesradio.se/podrss",
	Media:      "http://search.yahoo.com/mrss/",
	Version:    "2.0",
}
