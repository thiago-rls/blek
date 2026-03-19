package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	BaseURL     string `yaml:"base_url"`
	Author      string `yaml:"author"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{
		Title:       "My Blek Site",
		Description: "A simple static site built with Blek",
		BaseURL:     "http://localhost:8080",
		Author:      "Anonymous",
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}
