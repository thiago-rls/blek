package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSplitFrontmatter(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedFM   FrontMatter
		expectedBody string
	}{
		{
			name: "valid frontmatter",
			input: `---
title: "Hello World"
date: "2023-10-27"
---
This is the body.`,
			expectedFM: FrontMatter{
				Title: "Hello World",
				Date:  "2023-10-27",
			},
			expectedBody: "This is the body.",
		},
		{
			name:         "no frontmatter",
			input:        "Just some text without frontmatter.",
			expectedFM:   FrontMatter{},
			expectedBody: "Just some text without frontmatter.",
		},
		{
			name: "empty body",
			input: `---
title: "Empty Body"
date: "2023-10-27"
---`,
			expectedFM: FrontMatter{
				Title: "Empty Body",
				Date:  "2023-10-27",
			},
			expectedBody: "",
		},
		{
			name: "extra dashes in body",
			input: `---
title: "Dashes"
date: "2023-10-27"
---
Body with --- some dashes.`,
			expectedFM: FrontMatter{
				Title: "Dashes",
				Date:  "2023-10-27",
			},
			expectedBody: "Body with --- some dashes.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, body, err := splitFrontmatter([]byte(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(fm, tt.expectedFM) {
				t.Errorf("expected frontmatter %+v, got %+v", tt.expectedFM, fm)
			}
			if body != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, body)
			}
		})
	}
}

func TestValidateFrontmatter(t *testing.T) {
	tests := []struct {
		name    string
		fm      FrontMatter
		wantErr bool
	}{
		{
			name: "valid",
			fm: FrontMatter{
				Title: "Valid Title",
				Date:  "2023-10-27",
			},
			wantErr: false,
		},
		{
			name: "missing title",
			fm: FrontMatter{
				Date: "2023-10-27",
			},
			wantErr: true,
		},
		{
			name: "missing date",
			fm: FrontMatter{
				Title: "Valid Title",
			},
			wantErr: true,
		},
		{
			name:    "empty",
			fm:      FrontMatter{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFrontmatter(tt.fm, "test.md")
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFrontmatter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseFile_URL(t *testing.T) {
	// Create a temporary directory for test content
	tmpDir, err := os.MkdirTemp("", "blek-test-content")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		content  string
		slug     string
		section  string
		expected string
	}{
		{
			name: "posts section",
			content: `---
title: "My Post"
date: "2023-10-27"
---`,
			slug:     "my-post",
			section:  "posts",
			expected: "/posts/my-post/",
		},
		{
			name: "no section (top level)",
			content: `---
title: "About"
date: "2023-10-27"
---`,
			slug:     "about",
			section:  "",
			expected: "/about/",
		},
		{
			name: "custom section",
			content: `---
title: "Project A"
date: "2023-10-27"
---`,
			slug:     "project-a",
			section:  "projects",
			expected: "/projects/project-a/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.slug+".md")
			if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			page, err := ParseFile(path, tt.slug, tt.section)
			if err != nil {
				t.Fatalf("ParseFile() unexpected error: %v", err)
			}

			if page.URL != tt.expected {
				t.Errorf("expected URL %q, got %q", tt.expected, page.URL)
			}
		})
	}
}
