package mkpub

import (
	"errors"
	"net/url"
	"path/filepath"
	"slices"
	"strings"
)

type Status string

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
	StatusValidValues = []Status{
		StatusCompleted,
		StatusHiatus,
		StatusInactive,
		StatusOngoing,
	}
	ErrChapterBookNil                   = errors.New("chapter's book/parent does not exist")
	ErrChapterMissingPossibleIdentifier = errors.New("one of the following values must be defined: \"fileName\", \"title\", \"uniqueID\"")
	ErrContentParsedNil                 = errors.New("parsed content map is not initialized")
	ErrContentFormatNotFound            = errors.New("parsed content format does not exist")
)

func (s Status) Valid() bool {
	return slices.Contains(StatusValidValues, Status(strings.ToLower(string(s))))
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

// ExternalReference points to an external object, such as a hyperlink
// that directs the user to external donation pages.
type ExternalReference struct {
	Name        string
	Address     string
	IsHyperlink bool
}

// MediaAsset refers to assets such as images and videos that are included within a [Book]
type MediaAsset struct {
	Name          string
	AlternateText string
	Caption       string
}

// Profile can represent an individual author or contributor, or an
// organization such as a publisher.
type Profile struct {
	UniqueID         string
	Name             string
	NameAlternate    []string
	Roles            []string
	ShortDescription string
	About            string
	Images           []MediaAsset
	Links            []ExternalReference

	Parent  *Index
	Books   []*Book
	Content Content
}

func (p *Profile) InitDefaults(uniqueID string, parent *Index) {
	p.UniqueID = uniqueID
	p.Parent = parent
}

// SeriesInfo describes internal information of a series
type SeriesInfo struct {
	Name             string
	NameAlternate    []string
	ShortDescription string
	About            string
	IDs              map[string]string
	Links            []ExternalReference
	Images           []MediaAsset
}

// BookSeriesItem describes a [Book]'s relation/entry in a series.
type BookSeriesItem struct {
	SeriesInfo

	IndexID     string
	EntryNumber float64
}

// SeriesIndex describes a series: a set of related books.
type SeriesIndex struct {
	SeriesInfo

	UniqueID string
	Parent   *Index
	Books    []*Book
	Content  Content
}

func (s *SeriesIndex) InitDefaults(uniqueID string, parent *Index) {
	s.UniqueID = uniqueID
	s.Parent = parent
}

// Index is the main object that contains all [Book]s and
// [SeriesIndex]es.
type Index struct {
	Title            string
	TitleAlternate   []string
	ShortDescription string
	LanguageCode     string
	URL              string
	FaviconImageName string
	Params           map[string]any

	LayoutsDirectory string
	InputPath        string
	Books            []Book
	Series           []SeriesIndex
	Profiles         []Profile
}

func (i *Index) InitDefaults(inputPath string) error {
	i.InputPath = inputPath

	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return err
	}

	i.Title = filepath.Base(absInputPath)
	i.LayoutsDirectory = filepath.Join(inputPath, "layout")

	return nil
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
	Status             Status
	Links              []ExternalReference
	Mirrors            []ExternalReference
	Tags               []string
	Series             []BookSeriesItem
	CoverImage         MediaAsset
	Images             []MediaAsset
	Copyright          string
	License            string
	DatePublishedStart string
	DatePublishedEnd   string
	IDs                map[string]string
	Params             map[string]any

	CharacterEncoding string
	UniqueID          string

	Parent    *Index
	InputPath string
	Content   Content
	Chapters  []Chapter
}

func (b *Book) InitDefaults(inputPath string, parent *Index) error {
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

func (b *Book) ChaptersFlattened() []*Chapter {
	return chaptersFlattened(&b.Chapters)
}

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
	Params           map[string]any

	UniqueID  string
	Book      *Book
	InputPath string
	Content   Content
	Chapters  []Chapter
	Next      *Chapter
	Previous  *Chapter
}

func (c *Chapter) ChaptersFlattened() []*Chapter {
	return chaptersFlattened(&c.Chapters)
}

func chaptersFlattened(chapters *[]Chapter) []*Chapter {
	var flattened []*Chapter

	for i := range *chapters {
		next := &((*chapters)[i])

		var nested []*Chapter
		if len(next.Chapters) > 0 {
			nested = next.ChaptersFlattened()
		}

		flattened = append(flattened, next)
		flattened = append(flattened, nested...)
	}

	return flattened
}
