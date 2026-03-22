package sources

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadSources(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "sources.csv")

	content := []byte("id,name,url,category\n123,Test Source,https://valid.com/rss,Technology\n")
	if err := os.WriteFile(csvPath, content, 0o644); err != nil {
		t.Fatal(err)
	}

	sources, err := ReadSources(csvPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(sources) != 1 {
		t.Fatalf("Expected 1 source, got %d", len(sources))
	}

	s := sources[0]
	if s.ID != "123" {
		t.Errorf("Expected ID '123', got '%s'", s.ID)
	}
	if s.Name != "Test Source" {
		t.Errorf("Expected Name 'Test Source', got '%s'", s.Name)
	}
	if s.URL != "https://valid.com/rss" {
		t.Errorf("Expected URL 'https://valid.com/rss', got '%s'", s.URL)
	}
	if s.Category != "Technology" {
		t.Errorf("Expected Category 'Technology', got '%s'", s.Category)
	}
}

func TestValidateSource(t *testing.T) {
	tests := []struct {
		name     string
		source   Source
		expected bool
	}{
		{
			name: "Valid source",
			source: Source{
				ID:       "hacker-news",
				Name:     "Hacker News",
				URL:      "https://news.ycombinator.com/rss",
				Category: "Technology",
			},
			expected: true,
		},
		{
			name: "Invalid ID",
			source: Source{
				ID:       "hacker news",
				Name:     "Hacker News",
				URL:      "https://news.ycombinator.com/rss",
				Category: "Technology",
			},
			expected: false,
		},
		{
			name: "Invalid Name",
			source: Source{
				ID:       "hacker-news",
				Name:     "Hacker @ News",
				URL:      "https://news.ycombinator.com/rss",
				Category: "Technology",
			},
			expected: false,
		},
		{
			name: "Invalid URL",
			source: Source{
				ID:       "hacker-news",
				Name:     "Hacker News",
				URL:      "news.ycombinator.com/rss",
				Category: "Technology",
			},
			expected: false,
		},
		{
			name: "Invalid Category",
			source: Source{
				ID:       "hacker-news",
				Name:     "Hacker News",
				URL:      "news.ycombinator.com/rss",
				Category: "technology",
			},
			expected: false,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := ValidateSource(tc.source)
			if actual != tc.expected {
				t.Errorf("Test %v - '%s' FAIL: expected: %v, got: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
