package sources

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"www.github.com/maxbrt/colibri/internal/utils"
)

type Source struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Category string `json:"category"`
}

func ReadSources(filePath string) ([]Source, error) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %s | %s", filePath, err)
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading CSV %s | %s", filePath, err)
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV file %s has no data rows", filePath)
	}

	header := records[0]
	colIndex := make(map[string]int, len(header))
	for i, col := range header {
		colIndex[col] = i
	}

	requiredCols := []string{"id", "name", "url", "category"}
	for _, col := range requiredCols {
		if _, ok := colIndex[col]; !ok {
			return nil, fmt.Errorf("CSV file %s missing required column: %s", filePath, col)
		}
	}

	var sources []Source
	for _, row := range records[1:] {
		sources = append(sources, Source{
			ID:       row[colIndex["id"]],
			Name:     row[colIndex["name"]],
			URL:      row[colIndex["url"]],
			Category: row[colIndex["category"]],
		})
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
