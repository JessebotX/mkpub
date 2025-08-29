package config

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

const (
	// StatusCompleted indicates the work is fully published, even if it
	// may receive small modifications in the future.
	StatusCompleted = "completed"

	// StatusHiatus indicates the work currently being published is
	// temporarily not receiving any updates for an abnormal duration of
	// time, and is planned to resume in the future.
	StatusHiatus = "hiatus"

	// StatusInactive indicates the incomplete published work is currently
	// not being worked on and will not be receiving any updates for the
	// foreseeable future.
	StatusInactive = "inactive"

	// StatusOngoing indicates the currently-incomplete work is still
	// being worked on and will receive potentially major updates in a
	// timely manner.
	StatusOngoing = "ongoing"
)

var (
	// StatusValidValues stores a list of possible values for the Status
	// config field.
	StatusValidValues = []string{
		StatusCompleted,
		StatusHiatus,
		StatusInactive,
		StatusOngoing,
	}
	ErrBookEmptyUniqueID     = errors.New("book: missing/empty UniqueID. UniqueID must have at least 1 non-space character")
	ErrBookEmptyLanguageCode = errors.New("book: missing/empty LanguageCode. LanguageCode must have at least 1 non-space character")
)

type ErrUnrecognizedStatus struct {
	UnrecognizedValue string
}

func (e ErrUnrecognizedStatus) Error() string {
	return fmt.Sprintf("status: failed to recogniez \"%s\". Status must be one of the following (case-insensitive): %v", e.UnrecognizedValue, strings.Join(StatusValidValues, ", "))
}

// Book is a written work that contains one or more [Chapter]s, which are
// usually defined in a specific reading order. It may also be apart of zero or
// more [BookSeriesItem]'s.
type Book struct {
	Title              string
	TitleAlternate     []string
	Subtitle           string
	TitleSort          string
	LanguageCode       string
	ShortDescription   string
	About              string
	URL                string
	Authors            []Profile
	AuthorsSort        string
	Contributors       []Profile
	Publishers         []Profile
	Status             string
	Links              []ExternalReference
	Mirrors            []ExternalReference
	Tags               []string
	Series             []BookSeriesItem
	CoverImage         Asset
	Copyright          string
	License            string
	DatePublishedStart time.Time
	DatePublishedEnd   time.Time
	IDs                map[string]string
	Extra              map[string]any

	UniqueID       string
	Assets         []Asset
	Parent         *Index
	InputDirectory string
	Content        Content
	Chapters       []Chapter
}

func (b *Book) SetDefaults(inputPath string, parent *Index) error {
	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return err
	}
	b.InputDirectory = absInputPath

	b.UniqueID = strings.TrimSpace(b.UniqueID)
	if b.UniqueID == "" {
		b.UniqueID = filepath.Base(b.InputDirectory)
		if b.UniqueID == "" {
			return ErrBookEmptyUniqueID
		}
	}

	b.LanguageCode = strings.TrimSpace(b.LanguageCode)
	if b.LanguageCode == "" {
		return ErrBookEmptyLanguageCode
	}

	// trim trailing spaces and ignore case
	b.Status = strings.ToLower(strings.TrimSpace(b.Status))
	if b.Status == "" {
		b.Status = StatusCompleted
	}

	if !slices.Contains(StatusValidValues, strings.ToLower(b.Status)) {
		return ErrUnrecognizedStatus{b.Status}
	}

	if b.Parent == nil {
		b.Parent = parent
	}

	if parent != nil {
		b.LanguageCode = parent.LanguageCode

		if parent.URL != "" {
			b.URL, _ = url.JoinPath(parent.URL, "books", b.UniqueID)
		}
	}

	// trim trailing spaces
	for i := range b.Tags {
		b.Tags[i] = strings.TrimSpace(b.Tags[i])
	}

	return nil
}

func (b *Book) ChaptersFlattened() []*Chapter {
	return chaptersFlattened(&b.Chapters)
}

func (b *Book) SetDatePublishedStart(input string) error {
	t, err := parseDateTime(input)
	if err != nil {
		return err
	}

	b.DatePublishedEnd = t

	return nil
}

func (b *Book) SetDatePublishedEnd(input string) error {
	t, err := parseDateTime(input)
	if err != nil {
		return err
	}

	b.DatePublishedEnd = t

	return nil
}
