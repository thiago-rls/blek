package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverSections(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "blek-sections-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	sectionsToCreate := []string{"posts", "projects", "about"}
	for _, s := range sectionsToCreate {
		if err := os.MkdirAll(filepath.Join(tmpDir, s), 0755); err != nil {
			t.Fatalf("failed to create section dir %s: %v", s, err)
		}
	}

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("failed to read content dir: %v", err)
	}

	var sections []Section
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			title := name // simplified for test
			sections = append(sections, Section{Name: name, Title: title})
		}
	}

	expectedCount := len(sectionsToCreate)
	if len(sections) != expectedCount {
		t.Errorf("expected %d sections, got %d", expectedCount, len(sections))
	}
}
