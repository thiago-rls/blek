package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const Version = "0.1.1"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		handleInit(os.Args[2:])
	case "build":
		handleBuild()
	case "serve":
		handleServe()
	case "clean":
		handleClean()
	case "new":
		handleNew(os.Args[2:])
	case "version":
		fmt.Printf("blek v%s\n", Version)
	case "help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}
}

func handleBuild() {
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		fmt.Println("Warning: config.yaml not found, using defaults.")
	}

	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Building site...")
	if err := Build("content", "output", "templates", cfg); err != nil {
		fmt.Printf("Build failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Build complete!")
}

func handleServe() {
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		fmt.Println("Warning: config.yaml not found, using defaults.")
	}

	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Starting development server...")
	if err := Serve(cfg, "content", "output", "templates", "static"); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}

func handleClean() {
	fmt.Println("Cleaning output directory...")
	if err := os.RemoveAll("output"); err != nil {
		fmt.Printf("Clean failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Done.")
}

func handleNew(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: blek new [post|page] <title>")
		os.Exit(1)
	}

	kind := args[0]
	title := args[1]
	slug := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	date := time.Now().Format("2006-01-02")

	var targetDir string
	var fileName string
	switch kind {
	case "post":
		targetDir = "content/posts"
		fileName = slug + ".md"
	case "page":
		targetDir = filepath.Join("content", slug)
		fileName = "index.md"
	default:
		fmt.Printf("Unknown content type: %s (expected 'post' or 'page')\n", kind)
		os.Exit(1)
	}

	filePath := filepath.Join(targetDir, fileName)

	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("Error: file already exists: %s\n", filePath)
		os.Exit(1)
	}

	content := fmt.Sprintf(`---
title: "%s"
date: "%s"
---

Write your %s content here.
`, title, date, kind)

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created new %s: %s\n", kind, filePath)
}

func printHelp() {
	fmt.Println("Blek - A simple static site generator")
	fmt.Println("\nUsage:")
	fmt.Println("  blek <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  build      Generate the static site in the output directory")
	fmt.Println("  serve      Build and start a local development server with auto-reload")
	fmt.Println("  clean      Remove the output directory")
	fmt.Println("  new        Create a new post or page")
	fmt.Println("  init       Initialize a new blek project in the current or specified directory")
	fmt.Println("  version    Show version information")
	fmt.Println("  help       Show this help message")
	fmt.Println("\nExample:")
	fmt.Println("  blek new post \"Awesome Post\"")
	fmt.Println("  blek new page \"About\"")
	fmt.Println("  blek serve")
}
