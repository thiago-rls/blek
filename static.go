package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func copyStatic(staticDir, outputDir string) error {
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		return nil // no static folder, nothing to do
	}

	return filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Build the destination path by replacing the staticDir prefix with outputDir
		rel, err := filepath.Rel(staticDir, path)
		if err != nil {
			return fmt.Errorf("resolving relative path: %w", err)
		}

		dest := filepath.Join(outputDir, rel)

		if info.IsDir() {
			return os.MkdirAll(dest, 0755)
		}

		return copyFile(path, dest)
	})
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("opening %s: %w", src, err)
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("creating %s: %w", dest, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copying %s to %s: %w", src, dest, err)
	}

	fmt.Printf("copied: %s\n", dest)
	return nil
}
