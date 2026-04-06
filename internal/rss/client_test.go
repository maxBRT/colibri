package rss

import (
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mmcdole/gofeed"
	s "www.github.com/maxbrt/colibri/internal/sources"
)

func TestFetchAndParse(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		postsLen      int
		errorExpected bool
	}{
		{
			name:          "Valid Feed",
			body:          ValidFeed,
			postsLen:      4,
			errorExpected: false,
		},
		{
			name:          "Valid Feed With no GUID",
			body:          ValidFeedNoGUID,
			postsLen:      4,
			errorExpected: false,
		},
		{
			name:          "Invalid Feed should not fail but be empty",
			body:          MissingTitleFeed,
			postsLen:      0,
			errorExpected: false,
		},
		{
			name:          "Empty should fail",
			body:          "",
			postsLen:      0,
			errorExpected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := MockFeedServer(tc.body)
			defer server.Close()
			source := s.Source{
				ID:       "test",
				Name:     "Test",
				URL:      server.URL,
				Category: "Test",
			}
			posts, err := FetchAndParse(source)
			if err != nil && !tc.errorExpected {
				t.Fatalf("Expected no error, got %v", err)
			}
			if len(posts) != tc.postsLen {
				t.Errorf("Expected %d post, got %d", tc.postsLen, len(posts))
			}
		})
	}
}

func TestFetchAndParse_Gzip(t *testing.T) {
	tc := struct {
		name     string
		body     string
		postsLen int
	}{
		name:     "Valid Feed",
		body:     ValidFeed,
		postsLen: 4,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		gw := gzip.NewWriter(w)
		gw.Write([]byte(tc.body))
		gw.Close()
	}))
	defer server.Close()

	source := s.Source{
		ID:       "test",
		Name:     "Test",
		URL:      server.URL,
		Category: "Test",
	}
	posts, err := FetchAndParse(source)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(posts) != tc.postsLen {
		t.Errorf("Expected %d post, got %d", tc.postsLen, len(posts))
	}
}

func TestFetchFeedMetadata(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		status        int
		errorExpected bool
		validate      func(t *testing.T, feed *gofeed.Feed)
	}{
		{
			name:          "Valid Feed With Image and iTunes",
			body:          ValidFeedWithImages,
			status:        http.StatusOK,
			errorExpected: false,
			validate: func(t *testing.T, feed *gofeed.Feed) {
				if feed.Image == nil || feed.Image.URL != "https://example.com/logo.png" {
					t.Fatalf("Expected feed image url, got %+v", feed.Image)
				}
				if feed.ITunesExt == nil || feed.ITunesExt.Image != "https://example.com/itunes.jpg" {
					t.Fatalf("Expected iTunes image, got %+v", feed.ITunesExt)
				}
			},
		},
		{
			name:          "Bad Status",
			body:          ValidFeedWithImages,
			status:        http.StatusInternalServerError,
			errorExpected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/rss+xml")
				w.WriteHeader(tc.status)
				w.Write([]byte(tc.body))
			}))
			defer server.Close()

			source := s.Source{
				ID:       "test",
				Name:     "Test",
				URL:      server.URL,
				Category: "Test",
			}
			feed, err := FetchFeedMetadata(source)
			if err != nil && !tc.errorExpected {
				t.Fatalf("Expected no error, got %v", err)
			}
			if err == nil && tc.errorExpected {
				t.Fatalf("Expected error, got none")
			}
			if err == nil && tc.validate != nil {
				tc.validate(t, feed)
			}
		})
	}
}

func MockFeedServer(body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}))
}

const ValidFeed = `
<rss version="2.0">
<channel>
<title>xkcd.com</title>
<link>https://xkcd.com/</link>
<description>xkcd.com: A webcomic of romance and math humor.</description>
<language>en</language>
<item>
<title>Home Remedies</title>
<link>https://xkcd.com/3217/</link>
<description><img src="https://imgs.xkcd.com/comics/home_remedies.png" title="As always, you are permitted to call one person for guidance, but that person must be a grandparent." alt="As always, you are permitted to call one person for guidance, but that person must be a grandparent." /></description>
<pubDate>Mon, 09 Mar 2026 04:00:00 -0000</pubDate>
<guid>https://xkcd.com/3217/</guid>
</item>
<item>
<title>Bazookasaurus</title>
<link>https://xkcd.com/3216/</link>
<description><img src="https://imgs.xkcd.com/comics/bazookasaurus.png" title="In contrast to the deep booming sound associated with the cannon in pop culture depictions, recent studies show it actually made more of a 'toot toot!' noise." alt="In contrast to the deep booming sound associated with the cannon in pop culture depictions, recent studies show it actually made more of a 'toot toot!' noise." /></description>
<pubDate>Fri, 06 Mar 2026 05:00:00 -0000</pubDate>
<guid>https://xkcd.com/3216/</guid>
</item>
<item>
<title>Solar Warning</title>
<link>https://xkcd.com/3215/</link>
<description><img src="https://imgs.xkcd.com/comics/solar_warning.png" title="This replaces the previous solar activity watch, which was issued last month when the sun took off its sunglasses." alt="This replaces the previous solar activity watch, which was issued last month when the sun took off its sunglasses." /></description>
<pubDate>Wed, 04 Mar 2026 05:00:00 -0000</pubDate>
<guid>https://xkcd.com/3215/</guid>
</item>
<item>
<title>Electric Vehicles</title>
<link>https://xkcd.com/3214/</link>
<description><img src="https://imgs.xkcd.com/comics/electric_vehicles.png" title="Now that I've finally gotten an electric vehicle, I'm never going back to an acoustic one." alt="Now that I've finally gotten an electric vehicle, I'm never going back to an acoustic one." /></description>
<pubDate>Mon, 02 Mar 2026 05:00:00 -0000</pubDate>
<guid>https://xkcd.com/3214/</guid>
</item>
</channel>
</rss>
`

const ValidFeedNoGUID = `
<rss version="2.0">
<channel>
<title>xkcd.com</title>
<link>https://xkcd.com/</link>
<description>xkcd.com: A webcomic of romance and math humor.</description>
<language>en</language>
<item>
<title>Home Remedies</title>
<link>https://xkcd.com/3217/</link>
<description><img src="https://imgs.xkcd.com/comics/home_remedies.png" title="As always, you are permitted to call one person for guidance, but that person must be a grandparent." alt="As always, you are permitted to call one person for guidance, but that person must be a grandparent." /></description>
<pubDate>Mon, 09 Mar 2026 04:00:00 -0000</pubDate>
</item>
<item>
<title>Bazookasaurus</title>
<link>https://xkcd.com/3216/</link>
<description><img src="https://imgs.xkcd.com/comics/bazookasaurus.png" title="In contrast to the deep booming sound associated with the cannon in pop culture depictions, recent studies show it actually made more of a 'toot toot!' noise." alt="In contrast to the deep booming sound associated with the cannon in pop culture depictions, recent studies show it actually made more of a 'toot toot!' noise." /></description>
<pubDate>Fri, 06 Mar 2026 05:00:00 -0000</pubDate>
</item>
<item>
<title>Solar Warning</title>
<link>https://xkcd.com/3215/</link>
<description><img src="https://imgs.xkcd.com/comics/solar_warning.png" title="This replaces the previous solar activity watch, which was issued last month when the sun took off its sunglasses." alt="This replaces the previous solar activity watch, which was issued last month when the sun took off its sunglasses." /></description>
<pubDate>Wed, 04 Mar 2026 05:00:00 -0000</pubDate>
</item>
<item>
<title>Electric Vehicles</title>
<link>https://xkcd.com/3214/</link>
<description><img src="https://imgs.xkcd.com/comics/electric_vehicles.png" title="Now that I've finally gotten an electric vehicle, I'm never going back to an acoustic one." alt="Now that I've finally gotten an electric vehicle, I'm never going back to an acoustic one." /></description>
<pubDate>Mon, 02 Mar 2026 05:00:00 -0000</pubDate>
</item>
</channel>
</rss>
`

const MissingTitleFeed = `
<rss version="2.0">
<channel>
<title>xkcd.com</title>
<link>https://xkcd.com/</link>
<description>xkcd.com: A webcomic of romance and math humor.</description>
<language>en</language>
<item>
<link>https://xkcd.com/3217/</link>
<description><img src="https://imgs.xkcd.com/comics/home_remedies.png" title="As always, you are permitted to call one person for guidance, but that person must be a grandparent." alt="As always, you are permitted to call one person for guidance, but that person must be a grandparent." /></description>
<pubDate>Mon, 09 Mar 2026 04:00:00 -0000</pubDate>
<guid>https://xkcd.com/3217/</guid>
</item>
<item>
<link>https://xkcd.com/3216/</link>
<description><img src="https://imgs.xkcd.com/comics/bazookasaurus.png" title="In contrast to the deep booming sound associated with the cannon in pop culture depictions, recent studies show it actually made more of a 'toot toot!' noise." alt="In contrast to the deep booming sound associated with the cannon in pop culture depictions, recent studies show it actually made more of a 'toot toot!' noise." /></description>
<pubDate>Fri, 06 Mar 2026 05:00:00 -0000</pubDate>
<guid>https://xkcd.com/3216/</guid>
</item>
<item>
<link>https://xkcd.com/3215/</link>
<description><img src="https://imgs.xkcd.com/comics/solar_warning.png" title="This replaces the previous solar activity watch, which was issued last month when the sun took off its sunglasses." alt="This replaces the previous solar activity watch, which was issued last month when the sun took off its sunglasses." /></description>
<pubDate>Wed, 04 Mar 2026 05:00:00 -0000</pubDate>
<guid>https://xkcd.com/3215/</guid>
</item>
<item>
<link>https://xkcd.com/3214/</link>
<description><img src="https://imgs.xkcd.com/comics/electric_vehicles.png" title="Now that I've finally gotten an electric vehicle, I'm never going back to an acoustic one." alt="Now that I've finally gotten an electric vehicle, I'm never going back to an acoustic one." /></description>
<pubDate>Mon, 02 Mar 2026 05:00:00 -0000</pubDate>
<guid>https://xkcd.com/3214/</guid>
</item>
</channel>
</rss>
`

const ValidFeedWithImages = `
<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">
<channel>
<title>Example Feed</title>
<link>https://example.com/</link>
<description>Example feed with images</description>
<image>
<url>https://example.com/logo.png</url>
<title>Example Logo</title>
<link>https://example.com/</link>
</image>
<itunes:image href="https://example.com/itunes.jpg" />
</channel>
</rss>
`
