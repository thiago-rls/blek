package main

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Build is the main entry point for the site generation.
func Build(contentDir, outputDir, templatesDir string, cfg *Config) error {
	// Clean output directory before building
	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("cleaning output dir: %w", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	if err := copyStatic("static", outputDir); err != nil {
		return fmt.Errorf("copying static assets: %w", err)
	}

	tmpl, err := LoadTemplates(templatesDir)
	if err != nil {
		return fmt.Errorf("loading templates: %w", err)
	}

	// Walk content directory, one subfolder = one section
	entries, err := os.ReadDir(contentDir)
	if err != nil {
		return fmt.Errorf("reading content dir: %w", err)
	}

	var posts []TemplateData

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		section := entry.Name() // "posts", "about", "projects", etc.
		sectionPath := filepath.Join(contentDir, section)

		if section == "posts" {
			collected, err := buildPosts(sectionPath, outputDir, tmpl, cfg)
			if err != nil {
				return fmt.Errorf("building posts: %w", err)
			}
			posts = append(posts, collected...)
		} else {
			if err := buildPages(sectionPath, section, outputDir, tmpl, cfg); err != nil {
				return fmt.Errorf("building section %s: %w", section, err)
			}
		}
	}

	// Sort posts by date, newest first
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].DateStr > posts[j].DateStr
	})

	// Render index page
	if err := renderIndex(outputDir, tmpl, posts, cfg); err != nil {
		return fmt.Errorf("building index: %w", err)
	}

	// Render RSS feed
	if err := renderRSS(outputDir, posts, cfg); err != nil {
		return fmt.Errorf("building rss: %w", err)
	}

	return nil
}

func buildPosts(sectionPath, outputDir string, tmpl *Templates, cfg *Config) ([]TemplateData, error) {
	files, err := os.ReadDir(sectionPath)
	if err != nil {
		return nil, fmt.Errorf("reading posts dir: %w", err)
	}

	var posts []TemplateData

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
			continue
		}

		slug := strings.TrimSuffix(f.Name(), ".md")
		path := filepath.Join(sectionPath, f.Name())

		page, err := ParseFile(path, slug, "posts")
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", path, err)
		}

		htmlBody, err := RenderMarkdown(page.Body)
		if err != nil {
			return nil, fmt.Errorf("rendering %s: %w", path, err)
		}

		data := TemplateData{
			Title:    page.Title,
			DateStr:  page.Date.Format("2006-01-02"),
			Tags:     page.Tags,
			HTMLBody: template.HTML(htmlBody),
			URL:      page.URL,
			Config:   cfg,
		}

		outPath := filepath.Join(outputDir, "posts", slug, "index.html")
		if err := writeFile(outPath, func(w io.Writer) error {
			return tmpl.RenderPost(w, data)
		}); err != nil {
			return nil, fmt.Errorf("writing post %s: %w", slug, err)
		}

		fmt.Printf("built post: %s\n", outPath)
		posts = append(posts, data)
	}

	return posts, nil
}

func buildPages(sectionPath, section, outputDir string, tmpl *Templates, cfg *Config) error {
	files, err := os.ReadDir(sectionPath)
	if err != nil {
		return fmt.Errorf("reading section dir: %w", err)
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
			continue
		}

		slug := strings.TrimSuffix(f.Name(), ".md")
		path := filepath.Join(sectionPath, f.Name())

		page, err := ParseFile(path, slug, section)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}

		htmlBody, err := RenderMarkdown(page.Body)
		if err != nil {
			return fmt.Errorf("rendering %s: %w", path, err)
		}

		data := TemplateData{
			Title:    page.Title,
			DateStr:  page.Date.Format("2006-01-02"),
			Tags:     page.Tags,
			HTMLBody: template.HTML(htmlBody),
			URL:      page.URL,
			Config:   cfg,
		}

		// index.md inside a section folder -> output/section/index.html
		// any other file -> output/section/slug/index.html
		var outPath string
		if slug == "index" {
			outPath = filepath.Join(outputDir, section, "index.html")
		} else {
			outPath = filepath.Join(outputDir, section, slug, "index.html")
		}

		if err := writeFile(outPath, func(w io.Writer) error {
			return tmpl.RenderPage(w, data)
		}); err != nil {
			return fmt.Errorf("writing page %s: %w", slug, err)
		}

		fmt.Printf("built page: %s\n", outPath)
	}

	return nil
}

func renderIndex(outputDir string, tmpl *Templates, posts []TemplateData, cfg *Config) error {
	data := IndexData{Posts: posts, Config: cfg}
	outPath := filepath.Join(outputDir, "index.html")
	return writeFile(outPath, func(w io.Writer) error {
		return tmpl.RenderIndex(w, data)
	})
}
