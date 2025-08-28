package book

import (
	"errors"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrChapterMissingUniqueID = errors.New("chapter: UniqueID missing")
	ErrChapterMissingBook     = errors.New("chapter: reference to parent Book missing")
)

type Chapter struct {
	UniqueID string
	Title    string

	Text              Content
	Subtitle          string
	ShortDescription  string
	Chapters          []Chapter
	AuthorsNotePrefix Content
	AuthorsNoteSuffix Content
	DatePublished     time.Time
	DateModified      time.Time
	Authors           []Profile
	Contributors      []Profile
	Mirrors           []ExternalReference
	Copyright         string
	Licenses          []string
	Extra             map[string]any

	FileName  string
	Assets    []Asset
	InputPath string
	Book      *Book
	Next      *Chapter
	Previous  *Chapter
}

func (c *Chapter) New(inputPath string, uniqueID string, title string, book *Book) error {
	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return err
	}
	c.InputPath = absInputPath

	c.UniqueID = strings.TrimSpace(uniqueID)
	if c.UniqueID == "" {
		c.UniqueID = filepath.Base(c.InputPath)
	}

	c.Title = title
	c.Book = book
	if c.Book == nil {
		return ErrChapterMissingBook
	}

	return nil
}

func (c *Chapter) EnsureDefaults() error {
	absInputPath, err := filepath.Abs(c.InputPath)
	if err != nil {
		return err
	}
	c.InputPath = absInputPath

	if c.Book == nil {
		return ErrChapterMissingBook
	}

	c.UniqueID = strings.TrimSpace(c.UniqueID)

	if c.UniqueID == "" && c.Title == "" {
		return ErrChapterMissingUniqueID
	}

	if c.UniqueID == "" {
		c.UniqueID = strings.TrimSpace(c.Title)

		if c.UniqueID == "" {
			return ErrChapterMissingUniqueID
		}
	}

	if c.Title == "" {
		c.Title = c.UniqueID
	}

	return nil
}
