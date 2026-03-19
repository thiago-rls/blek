package main

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

type Section struct {
	Name  string
	Title string
}

type TemplateData struct {
	Title    string
	DateStr  string
	Tags     []string
	HTMLBody template.HTML
	URL      string
	Config   *Config
	Sections []Section
}

type IndexData struct {
	Posts    []TemplateData
	Config   *Config
	Sections []Section
}

type Templates struct {
	post  *template.Template
	page  *template.Template
	index *template.Template
}

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

func (t *Templates) RenderPost(w io.Writer, data TemplateData) error {
	return t.post.ExecuteTemplate(w, "base.html", data)
}

func (t *Templates) RenderPage(w io.Writer, data TemplateData) error {
	return t.page.ExecuteTemplate(w, "base.html", data)
}

func (t *Templates) RenderIndex(w io.Writer, data IndexData) error {
	return t.index.ExecuteTemplate(w, "base.html", data)
}

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
