package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed skeleton
var skeletonFS embed.FS

func handleInit(args []string) {
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	fmt.Printf("Initializing new blek project in %s...\n", targetDir)

	// Ensure target directory exists
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Printf("Error creating target directory: %v\n", err)
		os.Exit(1)
	}

	// We walk through the 'skeleton' directory in the embedded FS
	err := fs.WalkDir(skeletonFS, "skeleton", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate the destination path
		// Remove the 'skeleton/' prefix from the path
		relPath, _ := filepath.Rel("skeleton", path)
		if relPath == "." {
			return nil
		}

		destPath := filepath.Join(targetDir, relPath)

		if d.IsDir() {
			// Create the directory
			return os.MkdirAll(destPath, 0755)
		}

		// It's a file, copy it
		if _, err := os.Stat(destPath); err == nil {
			fmt.Printf("Skipping existing file: %s\n", destPath)
			return nil
		}

		content, err := skeletonFS.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, content, 0644)
	})

	if err != nil {
		fmt.Printf("Error initializing project: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Project initialized successfully!")
	fmt.Println("\nTo get started:")
	if targetDir != "." {
		fmt.Printf("  cd %s\n", targetDir)
	}
	fmt.Println("  blek build")
	fmt.Println("  blek serve")
}
