package pub

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

var (
	ErrBookMissingLanguageCode = errors.New("missing LanguageCode")
	ErrBookMissingUniqueID     = errors.New("missing UniqueID (value must contain at least 1 non-space character)")
	ErrBookMissingTitles       = errors.New("missing at least 1 title in Titles")
)

type ErrBookDuplicateTitles struct {
	Duplicate  string
	Index      int
	OtherIndex int
}

func (e ErrBookDuplicateTitles) Error() string {
	return fmt.Sprintf("duplicate title \"%s\" found in Titles in entry numbers %d and %d", e.Duplicate, e.Index, e.OtherIndex)
}

type ErrBookEmptyTag struct {
	Input string
	Index int
}

func (e ErrBookEmptyTag) Error() string {
	if e.Input == "" {
		return fmt.Sprintf("tag in Tags entry number %d is empty (value must contain at least 1 non-space character)", e.Index)
	}

	return fmt.Sprintf("tag \"%s\" provided is empty after trailing spaces are trimmed (value must contain at least 1 non-space character)", e.Input)
}

// Book represents a written work, which generally has an ordered list of 1 or more [Chapter]s.
type Book struct {
	UniqueID string `json:"unique_id"`

	// A list of titles of the book. The first title specified
	// is considered the primary title while the rest are considered
	// as alternate titles.
	//
	// USAGE
	//
	// titles:
	//   - "Main Title"        # required
	//   - "Alternate Title 1" # optional
	//   - "Alternate Title 2" # optional
	Titles             []string       `json:"titles"`
	Subtitle           string         `json:"subtitle"`
	Description        string         `json:"description"`
	Tagline            string         `json:"tagline"`
	Content            Content        `json:"content"`
	Authors            []Profile      `json:"authors"`
	Contributors       []Profile      `json:"contributors"`
	Publishers         []Profile      `json:"publishers"`
	Tags               []string       `json:"tags"`
	Status             Status         `json:"status"`
	Series             []Series       `json:"series"`
	Edition            string         `json:"edition"`
	URL                string         `json:"url"`
	LanguageCode       string         `json:"language_code"`
	DatePublishedStart *DateTime      `json:"date_published_start"`
	DatePublishedEnd   *DateTime      `json:"date_published_end"`
	LinksFunding       []Reference    `json:"links_funding"`
	LinksMirrors       []Reference    `json:"links_mirrors"`
	LinksOther         []Reference    `json:"links_other"`
	Assets             []Asset        `json:"assets"`
	IDs                map[string]any `json:"ids"`
	Copyright          Copyright      `json:"copyright"`
	Chapters           []Chapter      `json:"chapters"`
	Extra              map[string]any `json:"extra"`

	InputPath string
}

func (b Book) Title() string {
	return b.Titles[0]
}

func (b *Book) SetInputPath(inputPath string) error {
	absPath, err := filepath.Abs(inputPath)
	if err != nil {
		return err
	}
	b.InputPath = absPath

	return nil
}

func (b *Book) SetDatePublishedStartFromString(input string) error {
	t, err := dateFromString(input)
	if err != nil {
		return err
	}
	*b.DatePublishedStart = t

	return nil
}

func (b *Book) SetDatePublishedEndFromString(input string) error {
	t, err := dateFromString(input)
	if err != nil {
		return err
	}
	*b.DatePublishedEnd = t

	return nil
}

func (b Book) AlternateTitles() []string {
	if len(b.Titles) == 1 {
		return nil
	}

	var titles []string

	for i := 1; i < len(b.Titles); i++ {
		title := b.Titles[i]
		titles = append(titles, title)
	}

	return titles
}

func (b *Book) SetUniqueID(uniqueID string) {
	b.UniqueID = strings.ToLower(strings.TrimSpace(uniqueID))
}

func (b *Book) AddTag(tag string) error {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return ErrBookEmptyTag{Index: -1, Input: tag}
	}

	// dont add duplicates
	if slices.Contains(b.Tags, tag) {
		return nil
	}
	b.Tags = append(b.Tags, tag)

	return nil
}

// Trim trailing spaces/newlines and remove duplicates in book tags list
func (b *Book) NormalizeAllTags() error {
	// trim trailing spaces and count occurrences of each tag
	m := make(map[string]int)
	for i := range b.Tags {
		b.Tags[i] = strings.ToLower(strings.TrimSpace(b.Tags[i]))
		if b.Tags[i] == "" {
			return ErrBookEmptyTag{Index: i + 1, Input: ""}
		}

		m[b.Tags[i]] += 1
	}

	// remove duplicates
	var newTagsList []string
	for tag, amount := range m {
		if amount == 1 {
			newTagsList = append(newTagsList, tag)
		}
	}
	b.Tags = newTagsList

	return nil
}

func (b Book) ChaptersAndSubchapters() []*Chapter {
	return allChapters(&b.Chapters)
}

func (b *Book) EnsureValid() error {
	// trim unnecessary characters from unique ID
	b.SetUniqueID(b.UniqueID)

	// use uniqueID as a title if there are no titles
	if len(b.Titles) == 0 && b.UniqueID != "" {
		b.Titles = append(b.Titles, b.UniqueID)
	}

	if b.UniqueID == "" {
		return ErrBookMissingUniqueID
	}

	if len(b.Titles) == 0 {
		return ErrBookMissingTitles
	}

	if b.LanguageCode == "" {
		return ErrBookMissingLanguageCode
	}

	// error when duplicates are found
	for i, title := range b.Titles {
		for j := i + 1; j < len(b.Titles); j++ {
			other := b.Titles[j]
			if title == other {
				return ErrBookDuplicateTitles{Duplicate: title, Index: i + 1, OtherIndex: j + 1}
			}
		}
	}

	var profiles []Profile
	profiles = append(profiles, b.Authors...)
	profiles = append(profiles, b.Publishers...)
	profiles = append(profiles, b.Contributors...)
	for _, profile := range profiles {
		if err := profile.EnsureValid(); err != nil {
			return err
		}
	}

	for _, series := range b.Series {
		if err := series.EnsureValid(); err != nil {
			return err
		}
	}

	var references []Reference
	references = append(references, b.LinksFunding...)
	references = append(references, b.LinksMirrors...)
	references = append(references, b.LinksOther...)

	for _, reference := range references {
		if err := reference.EnsureValid(); err != nil {
			return err
		}
	}

	for _, asset := range b.Assets {
		if err := asset.EnsureValid(); err != nil {
			return err
		}
	}

	for _, chapter := range b.Chapters {
		if err := chapter.EnsureValid(); err != nil {
			return err
		}
	}

	return nil
}
