package client

import (
	"encoding/xml"
	"time"
)

type EpisodeListing struct {
	XMLName    xml.Name `xml:"sr"`
	Text       string   `xml:",chardata"`
	Copyright  string   `xml:"copyright"`
	Pagination struct {
		Text       string `xml:",chardata"`
		Page       string `xml:"page"`
		Size       string `xml:"size"`
		Totalhits  string `xml:"totalhits"`
		Totalpages string `xml:"totalpages"`
		Nextpage   string `xml:"nextpage"`
	} `xml:"pagination"`
	Episodes struct {
		Text    string    `xml:",chardata"`
		Episode []Episode `xml:"episode"`
	} `xml:"episodes"`
}

type Episode struct {
	Text        string `xml:",chardata"`
	ID          int    `xml:"id,attr"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
	URL         string `xml:"url"`
	Program     struct {
		Text string `xml:",chardata"`
		ID   int    `xml:"id,attr"`
		Name string `xml:"name,attr"`
	} `xml:"program"`
	Audiopreference   string    `xml:"audiopreference"`
	Audiopriority     string    `xml:"audiopriority"`
	Audiopresentation string    `xml:"audiopresentation"`
	Publishdateutc    time.Time `xml:"publishdateutc"`
	Imageurl          string    `xml:"imageurl"`
	Imageurltemplate  string    `xml:"imageurltemplate"`
	Listenpodfile     struct {
		Text            string    `xml:",chardata"`
		ID              string    `xml:"id,attr"`
		URL             string    `xml:"url"`
		Statkey         string    `xml:"statkey"`
		Duration        int       `xml:"duration"`
		Publishdateutc  time.Time `xml:"publishdateutc"`
		Title           string    `xml:"title"`
		Description     string    `xml:"description"`
		Filesizeinbytes string    `xml:"filesizeinbytes"`
		Program         struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"program"`
		Availablefromutc string `xml:"availablefromutc"`
	} `xml:"listenpodfile"`
	Downloadpodfile struct {
		Text            string    `xml:",chardata"`
		ID              string    `xml:"id,attr"`
		URL             string    `xml:"url"`
		Statkey         string    `xml:"statkey"`
		Duration        int       `xml:"duration"`
		Publishdateutc  time.Time `xml:"publishdateutc"`
		Title           string    `xml:"title"`
		Description     string    `xml:"description"`
		Filesizeinbytes int       `xml:"filesizeinbytes"`
		Program         struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"program"`
		Availablefromutc string `xml:"availablefromutc"`
	} `xml:"downloadpodfile"`
	Photographer string `xml:"photographer"`
	Broadcast    struct {
		Text             string `xml:",chardata"`
		Availablestoputc struct {
			Text string `xml:",chardata"`
			Nil  string `xml:"nil,attr"`
			P5   string `xml:"p5,attr"`
		} `xml:"availablestoputc"`
		Playlist struct {
			Text           string `xml:",chardata"`
			ID             string `xml:"id,attr"`
			URL            string `xml:"url"`
			Statkey        string `xml:"statkey"`
			Duration       string `xml:"duration"`
			Publishdateutc string `xml:"publishdateutc"`
		} `xml:"playlist"`
		Broadcastfiles struct {
			Text          string `xml:",chardata"`
			Broadcastfile struct {
				Text           string `xml:",chardata"`
				ID             string `xml:"id,attr"`
				URL            string `xml:"url"`
				Statkey        string `xml:"statkey"`
				Duration       string `xml:"duration"`
				Publishdateutc string `xml:"publishdateutc"`
			} `xml:"broadcastfile"`
		} `xml:"broadcastfiles"`
	} `xml:"broadcast"`
	Broadcasttime struct {
		Text         string `xml:",chardata"`
		Starttimeutc string `xml:"starttimeutc"`
		Endtimeutc   string `xml:"endtimeutc"`
	} `xml:"broadcasttime"`
}

type ProgramInfo struct {
	XMLName   xml.Name `xml:"sr"`
	Text      string   `xml:",chardata"`
	Copyright string   `xml:"copyright"`
	Program   struct {
		Text            string `xml:",chardata"`
		ID              int    `xml:"id,attr"`
		Name            string `xml:"name,attr"`
		Description     string `xml:"description"`
		Programcategory struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"programcategory"`
		Broadcastinfo            string `xml:"broadcastinfo"`
		Email                    string `xml:"email"`
		Phone                    string `xml:"phone"`
		Programurl               string `xml:"programurl"`
		Programslug              string `xml:"programslug"`
		Programimage             string `xml:"programimage"`
		Programimagetemplate     string `xml:"programimagetemplate"`
		Programimagewide         string `xml:"programimagewide"`
		Programimagetemplatewide string `xml:"programimagetemplatewide"`
		Socialimage              string `xml:"socialimage"`
		Socialimagetemplate      string `xml:"socialimagetemplate"`
		Socialmediaplatforms     struct {
			Text                string `xml:",chardata"`
			Socialmediaplatform []struct {
				Text        string `xml:",chardata"`
				Platform    string `xml:"platform,attr"`
				Platformurl string `xml:"platformurl,attr"`
			} `xml:"socialmediaplatform"`
		} `xml:"socialmediaplatforms"`
		Channel struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"channel"`
		Archived          string `xml:"archived"`
		Hasondemand       string `xml:"hasondemand"`
		Haspod            string `xml:"haspod"`
		Responsibleeditor string `xml:"responsibleeditor"`
	} `xml:"program"`
}
