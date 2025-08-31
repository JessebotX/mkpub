package config

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	ErrChapterBookNil                   = errors.New("chapter: book/parent does not exist")
	ErrChapterMissingPossibleIdentifier = errors.New("chapter: one of the following values must be defined: \"fileName\", \"title\", \"uniqueID\"")
)

// Chapter represents a division in a [Book].
type Chapter struct {
	Title            string
	TitleAlternate   []string
	Subtitle         string
	TitleSort        string
	ShortDescription string
	FileName         string
	AuthorsNote      string
	LanguageCode     string
	Authors          []Profile
	Contributors     []Profile
	AuthorsSort      string
	Links            []ExternalReference
	Mirrors          []ExternalReference
	DatePublished    string
	DateModified     string
	Copyright        string
	License          string
	IDs              map[string]string
	Extra            map[string]any

	UniqueID  string
	Book      *Book
	InputPath string
	Content   Content
	Chapters  []Chapter
	Next      *Chapter
	Previous  *Chapter
}

func (c *Chapter) SetDefaults(inputPath string, book *Book) error {
	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return err
	}
	c.InputPath = absInputPath

	c.UniqueID = strings.TrimSuffix(filepath.Base(c.InputPath), ".md")

	if c.LanguageCode == "" {
		c.LanguageCode = book.LanguageCode
	}
	c.Book = book

	return nil
}

func (c *Chapter) ChaptersFlattened() []*Chapter {
	return chaptersFlattened(&c.Chapters)
}

func (c *Chapter) HasSubchapters() bool {
	return len(c.Chapters) > 0
}
