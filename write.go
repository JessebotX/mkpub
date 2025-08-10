package mkpub

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func WriteIndexToStaticWebsite(index *OutputIndex, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// --- handle static files ---
	if err := copyDirectory(index.LayoutsDirectory, outputDir, []string{
		"index.html",
		"_book.html",
		"_chapter.html",
		"_author.html",
		"_series.html",
	}); err != nil {
		return fmt.Errorf("write index: failed to copy files to output: %w", err)
	}

	return nil
}

func copyDirectory(sourcePath, destinationPath string, excludePaths []string) error {
	return copyDirectoryHelper(sourcePath, destinationPath, sourcePath, excludePaths)
}

func copyDirectoryHelper(curr, dest, start string, excludePaths []string) error {
	items, err := os.ReadDir(curr)
	if err != nil {
		return err
	}

	for _, item := range items {
		target := filepath.Join(curr, item.Name())
		targetFromStart := strings.TrimLeft(strings.TrimPrefix(target, start), "/\\")

		// check exclusions
		if slices.Contains(excludePaths, targetFromStart) {
			continue
		}

		newFilePath := filepath.Join(dest, targetFromStart)

		// recursively copy subdirectories
		if item.IsDir() {
			if err := os.MkdirAll(newFilePath, 0755); err != nil {
				return err
			}

			if err := copyDirectoryHelper(target, dest, start, excludePaths); err != nil {
				return err
			}
		} else {
			if err := copyFile(target, newFilePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(sourcePath, destinationPath string) error {
	in, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return nil
}
