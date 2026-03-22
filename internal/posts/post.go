package posts

import (
	"fmt"
	"time"
)

type Post struct {
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description,omitempty"`
	PubDate     time.Time `json:"pubDate"`
	GUID        string    `json:"guid"`
	SourceID    string    `json:"sourceId"`
}

func NewPost(
	title string,
	link string,
	desc string,
	pubDate time.Time,
	guid string,
	sourceID string,
) *Post {
	if guid == "" {
		guid = fmt.Sprintf("%s.%s", link, "guid")
	}

	return &Post{
		Title:       title,
		Link:        link,
		Description: desc,
		PubDate:     pubDate,
		GUID:        guid,
		SourceID:    sourceID,
	}
}
