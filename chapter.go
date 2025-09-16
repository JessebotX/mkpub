package pub

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	ErrChapterMissingBookPointer = errors.New("missing pointer to Book")
	ErrChapterMissingUniqueID    = errors.New("missing UniqueID (must have at least 1 non-space character)")
)

// Chapter represents a division in a [Book] that contains its primary [Content].
type Chapter struct {
	UniqueID          string            `json:"unique_id"`
	Title             string            `json:"title"`
	Subtitle          string            `json:"subtitle"`
	Authors           []Profile         `json:"authors"`
	Contributors      []Profile         `json:"contributors"`
	Publishers        []Profile         `json:"publishers"`
	ContentFileName   string            `json:"content_file_name"`
	Content           Content           `json:"content"`
	AuthorsNotePrefix Content           `json:"authors_note_prefix"`
	AuthorsNoteSuffix Content           `json:"authors_note_suffix"`
	URL               string            `json:"url"`
	LanguageCode      string            `json:"language_code"`
	DatePublished     *DateTime         `json:"date_published"`
	DateUpdated       *DateTime         `json:"date_updated"`
	IDs               map[string]string `json:"ids"`
	Copyright         Copyright         `json:"copyright"`
	Extra             map[string]any    `json:"extra"`
	Chapters          []Chapter         `json:"chapters"`

	Previous  *Chapter
	Next      *Chapter
	Book      *Book
	InputPath string
}

func (c *Chapter) SetBook(book *Book) error {
	if book == nil {
		return ErrChapterMissingBookPointer
	}

	c.Book = book

	return nil
}

func (c Chapter) GetAbsoluteInputPath() (string, error) {
	absPath, err := filepath.Abs(c.InputPath)
	if err != nil {
		return c.InputPath, err
	}

	return absPath, nil
}

func (c *Chapter) SetInputPath(inputPath string) error {
	absPath, err := filepath.Abs(inputPath)
	if err != nil {
		return err
	}
	c.InputPath = absPath

	return nil
}

func (c *Chapter) SetUniqueID(uniqueID string) {
	c.UniqueID = strings.ToLower(strings.TrimSpace(uniqueID))
}

func (c *Chapter) SetDatePublishedFromString(input string) error {
	t, err := dateFromString(input)
	if err != nil {
		return err
	}
	*c.DatePublished = t

	return nil
}

func (c *Chapter) SetDateUpdatedFromString(input string) error {
	t, err := dateFromString(input)
	if err != nil {
		return err
	}
	*c.DateUpdated = t

	return nil
}

func (c Chapter) HasSubchapters() bool {
	return len(c.Chapters) > 0
}

func (c Chapter) Subchapters() []*Chapter {
	return allChapters(&c.Chapters)
}

func (c *Chapter) EnsureValid() error {
	if c.Book == nil {
		return ErrChapterMissingBookPointer
	}

	c.SetUniqueID(c.UniqueID)
	if c.UniqueID == "" && c.Title == "" && c.ContentFileName == "" {
		return ErrChapterMissingUniqueID
	}

	if c.UniqueID == "" && c.Title != "" {
		c.SetUniqueID(c.Title)
	}

	if c.UniqueID == "" && c.ContentFileName != "" {
		base := filepath.Base(c.ContentFileName)
		ext := filepath.Ext(base)

		c.SetUniqueID(strings.TrimSuffix(base, ext))
	}

	if c.Title == "" {
		c.Title = c.UniqueID
	}

	var profiles []Profile
	profiles = append(profiles, c.Authors...)
	profiles = append(profiles, c.Publishers...)
	profiles = append(profiles, c.Contributors...)
	for _, profile := range profiles {
		if err := profile.EnsureValid(); err != nil {
			return err
		}
	}

	for _, subchapter := range c.Chapters {
		if err := subchapter.EnsureValid(); err != nil {
			return err
		}
	}

	return nil
}
