// Package rss
package rss

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"www.github.com/maxbrt/colibri/internal/utils"
)

type Source struct {
	ID       string `yaml:"id"`
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Category string `yaml:"category"`
}

func ReadSources(dirPath string) ([]Source, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Printf("Error reading directory %s | %s", dirPath, err)
		return nil, err
	}

	var sources []Source

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		filePath := filepath.Join(dirPath, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading file %s | %s", filePath, err)
			return nil, err
		}

		s := Source{}
		if err := yaml.Unmarshal([]byte(content), &s); err != nil {
			log.Printf("Error unmarshalling file %s | %s", filePath, err)
			return nil, err
		}

		sources = append(sources, s)
	}

	return sources, nil
}

func ValidateSource(s Source) bool {
	IDIsValid := utils.IsValidSlug(s.ID)
	NameIsValid := utils.IsValidSourceName(s.Name)
	URLIsValid := utils.IsValidURL(s.URL)
	CategoryIsValid := utils.IsValidCategory(s.Category)

	return IDIsValid && NameIsValid && URLIsValid && CategoryIsValid
}
