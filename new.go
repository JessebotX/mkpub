package pub

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

const (
	BookConfigFileName         = "pub.yml"
	BookChaptersConfigFileName = "nav.yml"
	BookAssetsDirName          = "assets"
	BookChaptersDirName        = "chapters"
)

func NewBook(inputPath string) (Book, error) {
	var book Book
	if err := book.SetInputPath(inputPath); err != nil {
		return book, fmt.Errorf("[BOOK] \"%s\": %w", inputPath, err)
	}

	if err := unmarshalFromYAMLFile(filepath.Join(book.InputPath, BookConfigFileName), &book); err != nil {
		return book, fmt.Errorf("[BOOK] \"%s\": %w", inputPath, err)
	}

	assets, err := newAssets(filepath.Join(book.InputPath, BookAssetsDirName))
	if err != nil {
		return book, fmt.Errorf("[BOOK] \"%s\": %w", inputPath, err)
	}
	book.Assets = assets

	if len(book.Content.Raw) == 0 {
		raw, err := os.ReadFile(filepath.Join(book.InputPath, "index.md"))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return book, fmt.Errorf("[BOOK] \"%s\": %w", inputPath, err)
		}
		book.Content.Raw = raw
	}

	chapters, err := newChapters(book.InputPath, &book)
	if err != nil {
		return book, fmt.Errorf("[BOOK] \"%s\": %w", inputPath, err)
	}
	book.Chapters = chapters

	if err := book.EnsureValid(); err != nil {
		return book, fmt.Errorf("[BOOK] \"%s\": %w", inputPath, err)
	}

	return book, nil
}

func newAssets(inputPath string) ([]Asset, error) {
	var assets []Asset

	items, err := os.ReadDir(inputPath)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.IsDir() {
			continue
		}

		// TODO: properly create asset (detecting mime type format, etc.)
		asset := Asset{
			Objects: []AssetDescriptor{
				{
					Name: item.Name(),
				},
			},
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

func newChapters(booksDir string, book *Book) ([]Chapter, error) {
	var chapters []Chapter
	if err := unmarshalFromYAMLFile(filepath.Join(booksDir, BookChaptersConfigFileName), &chapters); err != nil {
		return chapters, err
	}

	var allChapters []*Chapter
	for i := range chapters {
		chapter := &chapters[i]
		if err := decodeChapter(chapter, book, filepath.Join(booksDir, BookChaptersDirName)); err != nil {
			return chapters, err
		}
		allChapters = append(allChapters, chapter)
		allChapters = append(allChapters, chapter.Subchapters()...)
	}

	for i := range allChapters {
		chapter := allChapters[i]

		if i-1 >= 0 {
			chapter.Previous = allChapters[i-1]
		}

		if i+1 < len(allChapters) {
			chapter.Next = allChapters[i+1]
		}
	}

	return chapters, nil
}

func decodeChapter(chapter *Chapter, book *Book, chaptersDir string) error {
	if err := chapter.SetBook(book); err != nil {
		return fmt.Errorf("[CHAPTER] \"%s\": %w", chapter.InputPath, err)
	}

	if chapter.InputPath == "" && chapter.ContentFileName != "" {
		if err := chapter.SetInputPath(filepath.Join(chaptersDir, chapter.ContentFileName)); err != nil {
			return fmt.Errorf("[CHAPTER] \"%s\": %w", chapter.ContentFileName, err)
		}
	}

	if chapter.InputPath != "" {
		raw, err := os.ReadFile(chapter.InputPath)
		if err != nil {
			return fmt.Errorf("[CHAPTER] \"%s\": %w", chapter.InputPath, err)
		}
		chapter.Content.Raw = raw
	}

	for i := range chapter.Subchapters() {
		subchapter := chapter.Subchapters()[i]
		if err := decodeChapter(subchapter, book, chaptersDir); err != nil {
			return err
		}
	}

	if err := chapter.EnsureValid(); err != nil {
		return fmt.Errorf("[CHAPTER] \"%s\": %w", chapter.InputPath, err)
	}

	return nil
}

func unmarshalFromYAMLFile(inputPath string, m any) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("parsing \"%s\": %w", filepath.Base(inputPath), err)
	}

	if err := yaml.Unmarshal(data, m); err != nil {
		return fmt.Errorf("parsing \"%s\": %w", filepath.Base(inputPath), err)
	}

	return nil
}
