package rss

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
	p "www.github.com/maxbrt/colibri/internal/posts"
	s "www.github.com/maxbrt/colibri/internal/sources"
	"www.github.com/maxbrt/colibri/internal/utils"
)

func FetchAndParse(source s.Source) ([]*p.Post, error) {
	c := NewRSSClient()
	req, err := NewRSSFetchRequest(source.URL)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		log.Printf("Failed to fetch feed %s | %s", source.ID, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	limitReader := io.LimitReader(resp.Body, 5*1024*1024)

	fp := gofeed.NewParser()
	feed, err := fp.Parse(limitReader)
	if err != nil {
		log.Printf("Failed to parse feed %s | %s", source.ID, err)
		return nil, err
	}

	var posts []*p.Post
	for _, i := range feed.Items {
		if !utils.IsValidFeedItem(*i) {
			continue
		}

		post := p.NewPost(
			i.Title,
			i.Link,
			"",
			*i.PublishedParsed,
			i.GUID,
			source.ID,
		)
		posts = append(posts, post)
	}

	return posts, nil
}

func NewRSSClient() *http.Client {
	return &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
}

func FetchFeedMetadata(source s.Source) (*gofeed.Feed, error) {
	c := NewRSSClient()
	req, err := NewRSSFetchRequest(source.URL)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		log.Printf("Failed to fetch feed %s | %s", source.ID, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	limitReader := io.LimitReader(resp.Body, 5*1024*1024)

	fp := gofeed.NewParser()
	feed, err := fp.Parse(limitReader)
	if err != nil {
		log.Printf("Failed to parse feed %s | %s", source.ID, err)
		return nil, err
	}

	return feed, nil
}

func NewRSSFetchRequest(URL string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", URL, nil)
	if err != nil {
		log.Printf("Error creating request: %s", err)
		return &http.Request{}, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MyService/1.0)")
	return req, nil
}
