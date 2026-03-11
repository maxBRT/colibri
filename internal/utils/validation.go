// Package utils regroup all utility functions for the project
package utils

import (
	"net/url"
	"regexp"
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
