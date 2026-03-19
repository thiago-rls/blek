package main

import (
	"os"
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Test when config.yaml is missing
	cfg, err := LoadConfig("non-existent.yaml")
	if err != nil {
		t.Fatalf("Expected no error when config is missing, got: %v", err)
	}

	if cfg.Title != "My Blek Site" {
		t.Errorf("Expected default title 'My Blek Site', got '%s'", cfg.Title)
	}
	if cfg.Author != "Anonymous" {
		t.Errorf("Expected default author 'Anonymous', got '%s'", cfg.Author)
	}
}

func TestLoadConfig_Partial(t *testing.T) {
	// Create a partial config file
	tmpFile := "partial_config.yaml"
	content := "title: \"Custom Title\"\n"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create partial config file: %v", err)
	}
	defer os.Remove(tmpFile)

	cfg, err := LoadConfig(tmpFile)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if cfg.Title != "Custom Title" {
		t.Errorf("Expected title 'Custom Title', got '%s'", cfg.Title)
	}
	if cfg.Author != "Anonymous" {
		t.Errorf("Expected default author 'Anonymous', got '%s'", cfg.Author)
	}
}
