package rss

import "time"

type Post struct {
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
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
	return &Post{
		Title:       title,
		Link:        link,
		Description: desc,
		PubDate:     pubDate,
		GUID:        guid,
		SourceID:    sourceID,
	}
}
