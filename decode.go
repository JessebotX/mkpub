package mkpub

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

const (
	IndexConfigName   = "mkpub.yml"
	BookConfigName    = "mkpub-book.yml"
	BookNavConfigName = "mkpub-book-nav.yml"
)

var (
	ErrChapterBookNil                   = errors.New("chapter's book/parent does not exist")
	ErrChapterMissingPossibleIdentifier = errors.New("one of the following values must be defined: \"fileName\", \"title\", \"uniqueID\"")
	ErrContentParsedNil                 = errors.New("parsed content map is not initialized")
	ErrContentFormatNotFound            = errors.New("parsed content format does not exist")
)

type OutputIndex struct {
	Index

	InputPath string
	Books     []OutputBook
}

func (i *OutputIndex) InitDefaults(inputPath string) error {
	i.InputPath = inputPath

	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return err
	}

	i.Title = filepath.Base(absInputPath)
	i.LayoutsDirectory = filepath.Join(inputPath, "layout")

	return nil
}

type OutputBook struct {
	Book

	Parent    *OutputIndex
	InputPath string
	Content   Content
	Chapters  []OutputChapter
}

func (b *OutputBook) InitDefaults(inputPath string, parent *OutputIndex) error {
	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return err
	}

	b.InputPath = absInputPath

	b.UniqueID = filepath.Base(b.InputPath)
	b.Title = b.UniqueID
	b.Parent = parent

	if parent != nil {
		b.LanguageCode = parent.LanguageCode

		if parent.URL != "" {
			b.URL, _ = url.JoinPath(parent.URL, "books", b.UniqueID)
		}
	}

	return nil
}

type OutputChapter struct {
	Chapter

	Book      *OutputBook
	InputPath string
	Content   Content
	Chapters  []OutputChapter
	Next      *OutputChapter
	Previous  *OutputChapter
}

func (c *OutputChapter) InitDefaults(inputPath string, book *OutputBook) error {
	if book == nil {
		return ErrChapterBookNil
	}

	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return err
	}

	c.InputPath = absInputPath
	c.UniqueID = filepath.Base(c.InputPath)
	c.Title = c.UniqueID
	c.LanguageCode = book.LanguageCode

	return nil
}

type OutputSeries struct {
	SeriesIndex

	Parent  *OutputIndex
	Books   []*Books
	Content Content
}

type Content struct {
	Raw []byte

	parsed map[string]any
}

func (c *Content) Format(format string) (any, error) {
	if c.parsed == nil {
		return "", ErrContentParsedNil
	}

	res, ok := c.parsed[format]
	if !ok {
		return "", ErrContentFormatNotFound
	}

	return res, nil
}

func (c *Content) AddFormat(format string, content any) {
	if c.parsed == nil {
		c.parsed = make(map[string]any, 1)
	}

	c.parsed[format] = content
}

func DecodeIndex(inputPath string) (OutputIndex, error) {
	var index OutputIndex
	if err := index.InitDefaults(inputPath); err != nil {
		return index, fmt.Errorf("index: failed on initialization: %w", err)
	}

	// --- Unmarshal config file ---
	confBody, err := os.ReadFile(filepath.Join(inputPath, IndexConfigName))
	if err != nil {
		return index, fmt.Errorf("index: failed to read %s: %w", IndexConfigName, err)
	}

	var confMap map[string]any
	if err := yaml.Unmarshal(confBody, &confMap); err != nil {
		return index, fmt.Errorf("index: failed to parse %s: %w", IndexConfigName, err)
	}

	if err := mapToStruct(confMap, &index); err != nil {
		return index, fmt.Errorf("index: failed to parse %s: %w", IndexConfigName, err)
	}

	index.Params = confMap

	booksDir := filepath.Join(inputPath, "books")
	dirs, err := os.ReadDir(booksDir)
	if err != nil {
		return index, fmt.Errorf("index: failed to read books directory: %w", err)
	}

	for i := range dirs {
		dir := dirs[i]

		if !dir.IsDir() {
			continue
		}

		book, err := DecodeBook(filepath.Join(booksDir, dir.Name()), &index)
		if err != nil {
			return index, err
		}

		index.Books = append(index.Books, book)
	}

	return index, nil
}

func DecodeBook(inputPath string, parent *OutputIndex) (OutputBook, error) {
	var book OutputBook
	if err := book.InitDefaults(inputPath, parent); err != nil {
		return book, fmt.Errorf("book \"%s\": failed on initialization: %w", filepath.Base(inputPath), err)
	}

	// --- Unmarshal config file ---
	confBody, err := os.ReadFile(filepath.Join(inputPath, BookConfigName))
	if err != nil {
		return book, fmt.Errorf("book \"%s\": failed to read %s: %w", book.UniqueID, BookConfigName, err)
	}

	var confMap map[string]any
	if err := yaml.Unmarshal(confBody, &confMap); err != nil {
		return book, fmt.Errorf("book \"%s\": failed to parse %s: %w", book.UniqueID, BookConfigName, err)
	}

	if err := mapToStruct(confMap, &book); err != nil {
		return book, fmt.Errorf("book \"%s\": failed to parse %s: %w", book.UniqueID, BookConfigName, err)
	}

	// --- Further parsing ---
	book.Params = confMap
	book.Content.Raw = []byte(book.About)

	if book.Status == "" {
		book.Status = StatusCompleted
	}

	if ok := book.Status.Valid(); !ok {
		return book, fmt.Errorf("book \"%s\": unrecognized status \"%s\". Must be one of the following (case-insensitive): %v", book.UniqueID, book.Status, StatusValidValues)
	}

	// --- Parse chapters ---
	navBody, err := os.ReadFile(filepath.Join(inputPath, BookNavConfigName))
	if err != nil {
		return book, fmt.Errorf("book \"%s\": failed to read %s: %w", book.UniqueID, BookNavConfigName, err)
	}

	var navConfMap []any
	if err := yaml.Unmarshal(navBody, &navConfMap); err != nil {
		return book, fmt.Errorf("book \"%s\": failed to parse %s: %w", book.UniqueID, BookNavConfigName, err)
	}

	var chapters []OutputChapter
	if err := mapToStruct(navConfMap, &chapters); err != nil {
		return book, fmt.Errorf("book \"%s\": failed to parse %s: %w", book.UniqueID, BookNavConfigName, err)
	}

	chaptersDir := filepath.Join(inputPath, "chapters")
	flattenedChapters, err := parseNav(&chapters, chaptersDir, &book)
	if err != nil {
		return book, fmt.Errorf("book \"%s\": %w", book.UniqueID, err)
	}

	// Set next and previous values
	for i := range flattenedChapters {
		c := flattenedChapters[i]

		if i-1 >= 0 {
			c.Previous = flattenedChapters[i-1]
		}

		if i+1 < len(flattenedChapters) {
			c.Next = flattenedChapters[i+1]
		}
	}

	book.Chapters = chapters

	return book, nil
}

// Returns a list of chapters in a flattened array for the purposes of pagination order.
func parseNav(chapters *[]OutputChapter, chaptersDir string, book *OutputBook) ([]*OutputChapter, error) {
	var flattenedList []*OutputChapter

	for i := range *chapters {
		c := &((*chapters)[i])
		c.InitDefaults(filepath.Join(chaptersDir, c.FileName), book)

		if c.FileName == "" && c.Title == "" && c.UniqueID == "" {
			return nil, ErrChapterMissingPossibleIdentifier
		}

		// If either title or uniqueID is missing...
		if c.Title != "" && c.UniqueID == "" {
			c.UniqueID = c.Title
		} else if c.Title == "" && c.UniqueID != "" {
			c.Title = c.UniqueID
		}

		if c.FileName != "" {
			raw, err := os.ReadFile(filepath.Join(chaptersDir, c.FileName))
			if err != nil {
				return nil, err
			}
			c.Content.Raw = raw
			c.InputPath = filepath.Join(chaptersDir, c.FileName)

			// if title or uniqueID is still missing, use fileName without the file extension
			if c.UniqueID == "" {
				c.UniqueID = strings.TrimSuffix(c.FileName, ".md")
			}

			if c.Title == "" {
				c.Title = strings.TrimSuffix(c.FileName, ".md")
			}
		}

		var nested []*OutputChapter
		if c.Chapters != nil {
			l, err := parseNav(&c.Chapters, chaptersDir, book)
			if err != nil {
				return nil, err
			}
			nested = l
		}

		flattenedList = append(flattenedList, c)
		flattenedList = append(flattenedList, nested...)
	}

	return flattenedList, nil
}

func mapToStruct[M map[string]any | []any](m M, s any) error {
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, s); err != nil {
		return err
	}

	return nil
}
