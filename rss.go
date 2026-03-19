package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func renderRSS(outputDir string, posts []TemplateData, cfg *Config) error {
	items := make([]Item, 0, len(posts))
	for _, p := range posts {
		items = append(items, Item{
			Title:       p.Title,
			Link:        cfg.BaseURL + p.URL,
			Description: p.Title,
			PubDate:     p.DateStr,
		})
	}

	feed := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:       cfg.Title,
			Link:        cfg.BaseURL,
			Description: cfg.Description,
			Items:       items,
		},
	}

	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling rss: %w", err)
	}

	outPath := filepath.Join(outputDir, "feed.xml")
	if err := os.WriteFile(outPath, append([]byte(xml.Header), data...), 0644); err != nil {
		return fmt.Errorf("writing rss: %w", err)
	}

	fmt.Printf("built rss: %s\n", outPath)
	return nil
}
