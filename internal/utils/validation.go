// Package utils regroup all utility functions for the project
package utils

import (
	"net/url"
	"regexp"

	"github.com/mmcdole/gofeed"
)

func IsValidSlug(s string) bool {
	regex := regexp.MustCompile("^[a-z0-9-]+$")
	return regex.MatchString(s)
}

func IsValidSourceName(s string) bool {
	regex := regexp.MustCompile("^[a-zA-Z0-9 ]+$")
	return regex.MatchString(s)
}

func IsValidURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Hostname() != "" && u.Path != ""
}

func IsValidCategory(s string) bool {
	regex := regexp.MustCompile("^[A-Z][a-z]*$")
	return regex.MatchString(s)
}

func IsValidFeedItem(i gofeed.Item) bool {
	title := i.Title
	link := i.Link
	pubDate := i.Published
	guid := i.GUID

	return title != "" && link != "" && pubDate != "" && guid != ""
}
