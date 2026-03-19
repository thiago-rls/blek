package main

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

// TemplateData is passed to post.html and page.html.
type TemplateData struct {
	Title    string
	DateStr  string
	Tags     []string
	HTMLBody template.HTML // Tells Go this string is safe, do not escape it
	URL      string
	Config   *Config
}

// IndexData is passed to index.html.
type IndexData struct {
	Posts  []TemplateData
	Config *Config
}

// Templates holds all loaded templates.
type Templates struct {
	post  *template.Template
	page  *template.Template
	index *template.Template
}

// LoadTemplates parses all templates from the templates/ directory.
func LoadTemplates(dir string) (*Templates, error) {
	base := filepath.Join(dir, "base.html")

	post, err := template.ParseFiles(base, filepath.Join(dir, "post.html"))
	if err != nil {
		return nil, fmt.Errorf("loading post template: %w", err)
	}

	page, err := template.ParseFiles(base, filepath.Join(dir, "page.html"))
	if err != nil {
		return nil, fmt.Errorf("loading page template: %w", err)
	}

	index, err := template.ParseFiles(base, filepath.Join(dir, "index.html"))
	if err != nil {
		return nil, fmt.Errorf("loading index template: %w", err)
	}

	return &Templates{post: post, page: page, index: index}, nil
}

// RenderPost writes a post page to w.
func (t *Templates) RenderPost(w io.Writer, data TemplateData) error {
	return t.post.ExecuteTemplate(w, "base.html", data)
}

// RenderPage writes a standalone page to w.
func (t *Templates) RenderPage(w io.Writer, data TemplateData) error {
	return t.page.ExecuteTemplate(w, "base.html", data)
}

// RenderIndex writes the index page to w.
func (t *Templates) RenderIndex(w io.Writer, data IndexData) error {
	return t.index.ExecuteTemplate(w, "base.html", data)
}

// writeFile creates the file at path and writes content to it.
func writeFile(path string, render func(io.Writer) error) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating directories for %s: %w", path, err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", path, err)
	}
	defer f.Close()
	return render(f)
}
