package mkpub

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/goccy/go-yaml"

	"github.com/JessebotX/mkpub/book"

	"golang.org/x/sync/errgroup"
)

const (
	MainIndexConfigName    = "mkpub.yml"
	MainIndexBooksDir      = "books"
	BookConfigName         = "book.yml"
	BookAssetsDir          = "assets"
	BookChaptersDir        = "chapters"
	BookChaptersConfigName = "nav.yml"
)

func DecodeMainIndex(inputPath string) (MainIndex, error) {
	var index MainIndex

	/*** Parse config ***/

	conf, err := yamlFileToMap(filepath.Join(inputPath, MainIndexConfigName))
	if err != nil {
		return index, fmt.Errorf("index: parse %s: %w", MainIndexConfigName, err)
	}
	defer clear(conf)

	if err := mapToStruct(conf, &index); err != nil {
		return index, fmt.Errorf("index: parse %s: %w", MainIndexConfigName, err)
	}

	/*** Parse books ***/

	books, err := bookDirsToBooks(filepath.Join(inputPath, MainIndexBooksDir), &index)
	if err != nil {
		return index, err
	}
	index.Books = books

	return index, nil
}

func bookDirsToBooks(booksDir string, index *MainIndex) ([]BookIndex, error) {
	var books []BookIndex

	items, err := os.ReadDir(booksDir)
	if err != nil {
		return books, fmt.Errorf("index: books directory: %w", err)
	}

	for _, item := range items {
		if !item.IsDir() {
			continue
		}

		var g errgroup.Group
		var book BookIndex

		g.Go(func() error {
			book, err = DecodeBook(filepath.Join(booksDir, item.Name()), index)
			if err != nil {
				return err
			}
			return nil
		})
		if err := g.Wait(); err != nil {
			return books, err
		}

		books = append(books, book)
	}

	return books, nil
}

func DecodeBook(inputPath string, index *MainIndex) (BookIndex, error) {
	b := BookIndex{
		Book: book.Book{
			InputPath: inputPath,
		},
		Index: index,
	}

	/*** Parse config ***/

	conf, err := yamlFileToMap(filepath.Join(inputPath, BookConfigName))
	if err != nil {
		return b, fmt.Errorf("book: parse %s: %w", BookConfigName, err)
	}
	defer clear(conf)

	f, err := os.ReadFile(filepath.Join(inputPath, BookConfigName))
	if err != nil {
		return b, err
	}

	if err := yaml.Unmarshal(f, &conf); err != nil {
		return b, fmt.Errorf("book: parse %s: %w", BookConfigName, err)
	}

	/*** Parse date strings ***/

	dateStartInput, ok := conf["publishedStart"]
	if ok {
		switch v := dateStartInput.(type) {
		case string:
			if err := b.ParseDatePublishedStart(v); err != nil {
				return b, err
			}
		case time.Time:
			b.DatePublishedStart = v
		default:
			return b, fmt.Errorf("book: parse %s: unsupported publishedStart type. Must be either a time.Time or a string", BookConfigName)
		}
	}

	dateEndInput, ok := conf["publishedEnd"]
	if ok {
		switch v := dateEndInput.(type) {
		case string:
			if err := b.ParseDatePublishedEnd(v); err != nil {
				return b, err
			}
		case time.Time:
			b.DatePublishedEnd = v
		default:
			return b, fmt.Errorf("book: parse %s: unsupported publishedEnd type. Must be either a time.Time or a string", BookConfigName)
		}
	}

	/*** Set defaults and ensure book is valid ***/

	if err := b.EnsureDefaults(); err != nil {
		return b, err
	}

	/*** parse chapters ***/
	chapters, err := decodeAllChapters(filepath.Join(inputPath, BookChaptersConfigName), filepath.Join(inputPath, BookChaptersDir), &b)
	if err != nil {
		return b, err
	}
	b.Chapters = chapters

	return b, nil
}

func decodeAllChapters(navFilePath, chaptersDir string, b *BookIndex) ([]book.Chapter, error) {
	var chapters []book.Chapter

	conf, err := yamlFileToArray(navFilePath)
	if err != nil {
		return chapters, fmt.Errorf("book: parse %s: %w", BookChaptersConfigName, err)
	}

	if err := mapToStruct(conf, &chapters); err != nil {
		return chapters, fmt.Errorf("book: parse %s: %w", BookChaptersConfigName, err)
	}

	var g errgroup.Group

	for i := range chapters {
		g.Go(func() error {
			return parseChapter(&chapters[i], chaptersDir, b)
		})
	}
	if err := g.Wait(); err != nil {
		return chapters, err
	}

	return chapters, nil
}

func parseChapter(c *book.Chapter, chaptersDir string, b *BookIndex) error {
	c.Book = &b.Book

	if c.FileName != "" {
		contents, err := os.ReadFile(filepath.Join(chaptersDir, c.FileName))
		if err != nil {
			return err
		}

		c.Text.Raw = contents
	}
	c.UniqueID = strings.TrimSuffix(c.FileName, ".md")
	if c.UniqueID == "" {
		c.UniqueID = strings.ToLower(c.Title)
	}

	if err := c.EnsureDefaults(); err != nil {
		return err
	}

	return nil
}

func yamlFileToArray(configPath string) ([]any, error) {
	var conf []any

	f, err := os.ReadFile(configPath)
	if err != nil {
		return conf, err
	}

	if err := yaml.Unmarshal(f, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}

func yamlFileToMap(configPath string) (map[string]any, error) {
	var conf map[string]any

	f, err := os.ReadFile(configPath)
	if err != nil {
		return conf, err
	}

	if err := yaml.Unmarshal(f, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}

func mapToStruct[T []any | map[string]any](m T, s any) error {
	jsonBody, err := json.Marshal(m)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonBody, s); err != nil {
		return err
	}

	return nil
}
