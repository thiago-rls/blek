package main

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRenderRSS(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "blek-rss-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &Config{
		Title:       "Test Blog",
		Description: "A test blog",
		BaseURL:     "https://example.com",
	}

	date, _ := time.Parse("2006-01-02", "2023-10-27")
	posts := []TemplateData{
		{
			Title: "Test Post",
			URL:   "/posts/test-post/",
			Date:  date,
		},
	}

	err = renderRSS(tmpDir, posts, cfg)
	if err != nil {
		t.Fatalf("renderRSS() unexpected error: %v", err)
	}

	rssPath := filepath.Join(tmpDir, "feed.xml")
	data, err := os.ReadFile(rssPath)
	if err != nil {
		t.Fatalf("failed to read generated rss: %v", err)
	}

	var rss RSS
	if err := xml.Unmarshal(data, &rss); err != nil {
		t.Fatalf("failed to unmarshal rss: %v", err)
	}

	if rss.Channel.Title != cfg.Title {
		t.Errorf("expected title %q, got %q", cfg.Title, rss.Channel.Title)
	}

	if len(rss.Channel.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(rss.Channel.Items))
	}

	expectedLink := "https://example.com/posts/test-post/"
	if rss.Channel.Items[0].Link != expectedLink {
		t.Errorf("expected link %q, got %q", expectedLink, rss.Channel.Items[0].Link)
	}
}
