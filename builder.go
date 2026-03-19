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

func Build(contentDir, outputDir, templatesDir string, cfg *Config) error {
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

	entries, err := os.ReadDir(contentDir)
	if err != nil {
		return fmt.Errorf("reading content dir: %w", err)
	}

	var sections []Section
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			title := strings.ToUpper(name[:1]) + name[1:]
			sections = append(sections, Section{Name: name, Title: title})
		}
	}

	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Title < sections[j].Title
	})

	var posts []TemplateData

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		section := entry.Name()
		sectionPath := filepath.Join(contentDir, section)

		if section == "posts" {
			collected, err := buildPosts(sectionPath, outputDir, tmpl, cfg, sections)
			if err != nil {
				return fmt.Errorf("building posts: %w", err)
			}
			posts = append(posts, collected...)
		} else {
			if err := buildPages(sectionPath, section, outputDir, tmpl, cfg, sections); err != nil {
				return fmt.Errorf("building section %s: %w", section, err)
			}
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].DateStr > posts[j].DateStr
	})

	if err := renderIndex(outputDir, tmpl, posts, cfg, sections); err != nil {
		return fmt.Errorf("building index: %w", err)
	}

	postsIndexDir := filepath.Join(outputDir, "posts")
	if err := renderIndex(postsIndexDir, tmpl, posts, cfg, sections); err != nil {
		return fmt.Errorf("building posts index: %w", err)
	}

	if err := renderRSS(outputDir, posts, cfg); err != nil {
		return fmt.Errorf("building rss: %w", err)
	}

	return nil
}

func buildPosts(sectionPath, outputDir string, tmpl *Templates, cfg *Config, sections []Section) ([]TemplateData, error) {
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
			Sections: sections,
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

func buildPages(sectionPath, section, outputDir string, tmpl *Templates, cfg *Config, sections []Section) error {
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
			Sections: sections,
		}

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

func renderIndex(outputDir string, tmpl *Templates, posts []TemplateData, cfg *Config, sections []Section) error {
	data := IndexData{Posts: posts, Config: cfg, Sections: sections}
	outPath := filepath.Join(outputDir, "index.html")
	return writeFile(outPath, func(w io.Writer) error {
		return tmpl.RenderIndex(w, data)
	})
}
