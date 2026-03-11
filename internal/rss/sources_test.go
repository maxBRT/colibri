package rss

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadSources(t *testing.T) {
	// Test that the read sources ignore everything but yaml files
	tmpDir := t.TempDir()

	yamlFilePath := filepath.Join(tmpDir, "test-source.yaml")
	content := []byte("id: 123\nname: Test Source\nurl: https://valid.com/rss\ncategory: Technology")
	if err := os.WriteFile(yamlFilePath, content, 0o644); err != nil {
		t.Fatal(err)
	}
	jsonFilePath := filepath.Join(tmpDir, "test-source.json")
	content = []byte("doesnt matter")
	if err := os.WriteFile(jsonFilePath, content, 0o644); err != nil {
		t.Fatal(err)
	}
	dirPath := filepath.Join(tmpDir, "directory")
	if err := os.Mkdir(dirPath, 0o755); err != nil {
		t.Fatal(err)
	}

	sources, err := ReadSources(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(sources) != 1 {
		t.Errorf("Expected 1 source, got %d", len(sources))
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
