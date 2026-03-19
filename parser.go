package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type FrontMatter struct {
	Title string   `yaml:"title"`
	Date  string   `yaml:"date"`
	Tags  []string `yaml:"tags"`
}

type Page struct {
	FrontMatter
	Body    string // raw markdown, not yet converted to HTML
	Slug    string
	URL     string
	Section string // "posts", "about", "projects", etc.
	Date    time.Time
}

func ParseFile(path, slug, section string) (*Page, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	fm, body, err := splitFrontmatter(data)
	if err != nil {
		return nil, fmt.Errorf("parsing frontmatter in %s: %w", path, err)
	}

	if err := validateFrontmatter(fm, path); err != nil {
		return nil, err
	}

	var parsed time.Time
	if fm.Date != "" {
		parsed, err = time.Parse("2006-01-02", fm.Date)
		if err != nil {
			return nil, fmt.Errorf("parsing date in %s: %w", path, err)
		}
	}

	url := "/" + section + "/" + slug + "/"
	if section == "" {
		url = "/" + slug + "/"
	}

	return &Page{
		FrontMatter: fm,
		Body:        body,
		Slug:        slug,
		URL:         url,
		Section:     section,
		Date:        parsed,
	}, nil
}

// splitFrontmatter separates the YAML block from the markdown body.
// It expects the file to start with ---, followed by YAML, closed by another ---.
func splitFrontmatter(data []byte) (FrontMatter, string, error) {
	var fm FrontMatter

	scanner := bufio.NewScanner(bytes.NewReader(data))

	scanner.Scan()
	if strings.TrimSpace(scanner.Text()) != "---" {
		return fm, string(data), nil
	}

	var yamlLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			break
		}

		yamlLines = append(yamlLines, line)
	}

	var bodyLines []string
	for scanner.Scan() {
		bodyLines = append(bodyLines, scanner.Text())
	}

	if err := yaml.Unmarshal([]byte(strings.Join(yamlLines, "\n")), &fm); err != nil {
		return fm, "", err
	}

	body := strings.TrimSpace(strings.Join(bodyLines, "\n"))
	return fm, body, nil
}

func validateFrontmatter(fm FrontMatter, path string) error {
	if fm.Title == "" {
		return fmt.Errorf("%s: missing title", path)
	}
	if fm.Date == "" {
		return fmt.Errorf("%s: missing date", path)
	}
	return nil
}
