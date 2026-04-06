package logo

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/thanhpk/go-favicon"
	ps "www.github.com/maxbrt/colibri/internal/pubsub"
	"www.github.com/maxbrt/colibri/internal/rss"
	s "www.github.com/maxbrt/colibri/internal/sources"
)

type logoCandidate struct {
	url      string
	mimeType string
}

func Listener(svc *Service) func(s.Source) ps.AckType {
	return func(src s.Source) ps.AckType {
		URL, err := url.Parse(src.URL)
		if err != nil {
			log.Printf("invalid url for source %s: %s", src.ID, src.URL)
			return ps.Ack
		}
		c := rss.NewRSSClient()

		candidate := logoCandidate{}
		feed, err := rss.FetchFeedMetadata(src)
		if err != nil {
			log.Printf("Failed to fetch feed metadata for %s: %s", src.ID, err)
		} else if feed != nil {
			if feed.Image != nil && feed.Image.URL != "" {
				candidate.url = feed.Image.URL
			} else if feed.ITunesExt != nil && feed.ITunesExt.Image != "" {
				candidate.url = feed.ITunesExt.Image
			}
		}

		if candidate.url == "" {
			fullURL := fmt.Sprintf("%s://%s", URL.Scheme, URL.Host)
			icons, err := favicon.Find(fullURL)
			if err != nil {
				log.Printf("Failed to find logo for %s: %s", fullURL, err)
				return ps.Ack
			}
			if len(icons) == 0 {
				log.Printf("No icons found for %s", fullURL)
				return ps.Ack
			}
			candidate.url = icons[0].URL
			candidate.mimeType = icons[0].MimeType
		}

		req, err := rss.NewRSSFetchRequest(candidate.url)
		if err != nil {
			log.Println(err)
			return ps.Ack
		}

		if candidate.mimeType != "" {
			req.Header.Add("Accept", fmt.Sprintf("image/%s", candidate.mimeType))
		}
		resp, err := c.Do(req)
		if err != nil {
			log.Printf("Failed to fetch logo for %s: %s", src.ID, err)
			if resp != nil && resp.StatusCode > http.StatusInternalServerError {
				return ps.NackRequeue
			}
			return ps.Ack
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Printf("Received %d when fetching logo for %s", resp.StatusCode, src.ID)
			if resp.StatusCode >= http.StatusInternalServerError {
				return ps.NackRequeue
			}
			return ps.Ack
		}
		limitReader := io.LimitReader(resp.Body, 5*1024*1024)
		b, err := io.ReadAll(limitReader)
		if err != nil {
			log.Printf("Failed to read logo form body: %s", err)
			return ps.Ack
		}
		contentType := candidate.mimeType
		if contentType == "" {
			contentType = resp.Header.Get("Content-Type")
		}
		detectedContentType := http.DetectContentType(b)
		if contentType == "" {
			contentType = detectedContentType
		}
		if contentType == "" || detectedContentType == "" || !isImageContentType(contentType, detectedContentType) {
			log.Printf(
				"Non-image response when fetching logo for %s (status=%d content-type=%q detected=%q body=%q)",
				src.ID,
				resp.StatusCode,
				contentType,
				detectedContentType,
				snippetBody(b),
			)
			return ps.Ack
		}
		_, err = svc.SaveLogo(context.Background(), src.ID, contentType, b)
		if err != nil {
			log.Printf("Failed to save logo for %s: %s", src.ID, err)
			return ps.NackRequeue
		}
		log.Printf("Stored logo for source %s", src.ID)

		return ps.Ack
	}
}

func isImageContentType(contentType string, detectedContentType string) bool {
	if len(contentType) >= 6 && contentType[:6] == "image/" {
		return true
	}
	if len(detectedContentType) >= 6 && detectedContentType[:6] == "image/" {
		return true
	}
	return false
}

func snippetBody(body []byte) string {
	const maxBytes = 512
	if len(body) <= maxBytes {
		return string(body)
	}
	return string(body[:maxBytes])
}
