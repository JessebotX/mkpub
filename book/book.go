package book

import (
	"errors"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

const (
	StatusOngoing   = "Ongoing"
	StatusInactive  = "Inactive"
	StatusHiatus    = "Hiatus"
	StatusCompleted = "Completed"
)

type ErrBookUnrecognizedStatus struct {
	UnrecognizedValue string
}

func (e ErrBookUnrecognizedStatus) Error() string {
	var msg strings.Builder

	_, _ = msg.WriteString("book: Status unrecognized: got = \"")
	_, _ = msg.WriteString(e.UnrecognizedValue)
	_, _ = msg.WriteString("\"; want = one of the following (case-insensitive): ")
	_, _ = msg.WriteString(strings.Join(StatusValidValues, ", "))

	return msg.String()
}

var (
	ErrBookMissingUniqueID     = errors.New("book: UniqueID missing")
	ErrBookEmptyUniqueID       = errors.New("book: UniqueID cannot be empty (must have at least 1 non-space character)")
	ErrBookMissingLanguageCode = errors.New("book: LanguageCodes missing. Want at least 1 language code")
	ErrBookEmptyLanguageCode   = errors.New("book: LanguageCodes cannot have empty values (must have at least 1 non-space character)")

	StatusValidValues = []string{
		strings.ToLower(StatusOngoing),
		strings.ToLower(StatusInactive),
		strings.ToLower(StatusHiatus),
		strings.ToLower(StatusCompleted),
	}
)

type Book struct {
	UniqueID      string
	Title         string
	LanguageCodes []string

	ShortDescription   string
	About              Content
	Tagline            string
	AlternateTitles    []string
	Subtitle           string
	TitleSort          string
	Authors            []Profile
	Contributors       []Profile
	Publishers         []Profile
	Tags               []string
	Series             []Series
	Edition            string
	DatePublishedStart time.Time
	DatePublishedEnd   time.Time
	IDs                map[string]string
	Status             string
	CoverImage         Asset
	DonationOptions    []ExternalReference
	Mirrors            []ExternalReference
	SocialLinks        []ExternalReference
	Copyright          string
	Licenses           []string
	Chapters           []Chapter
	Extra              map[string]any

	InputPath string
	Assets    []Asset
}

func (b *Book) New(inputPath string, uniqueID, title, languageCode string) error {
	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return err
	}
	b.InputPath = absInputPath

	b.UniqueID = strings.TrimSpace(uniqueID)
	if b.UniqueID == "" {
		b.UniqueID = filepath.Base(b.InputPath)
	}
	b.Title = strings.TrimSpace(title)
	b.LanguageCodes = append(b.LanguageCodes, languageCode)

	return nil
}

func (b *Book) EnsureDefaults() error {
	absInputPath, err := filepath.Abs(b.InputPath)
	if err != nil {
		return err
	}
	b.InputPath = absInputPath
	if b.UniqueID == "" {
		return ErrBookMissingUniqueID
	}

	b.UniqueID = strings.TrimSpace(b.UniqueID)
	if b.UniqueID == "" {
		return ErrBookMissingUniqueID
	}

	if len(b.LanguageCodes) == 0 {
		return ErrBookMissingLanguageCode
	}

	for i := range b.LanguageCodes {
		b.LanguageCodes[i] = strings.TrimSpace(b.LanguageCodes[i])
		if b.LanguageCodes[i] == "" {
			return ErrBookEmptyLanguageCode
		}
	}

	b.Status = strings.ToLower(strings.TrimSpace(b.Status))
	if !slices.Contains(StatusValidValues, b.Status) {
		return ErrBookUnrecognizedStatus{UnrecognizedValue: b.Status}
	}

	for i := range b.Tags {
		b.Tags[i] = strings.TrimSpace(b.Tags[i])
	}

	return nil
}

func (b *Book) ParseDatePublishedStart(input string) error {
	t, err := parseDateInput(input)
	if err != nil {
		return err
	}

	b.DatePublishedStart = t

	return nil
}

func (b *Book) ParseDatePublishedEnd(input string) error {
	t, err := parseDateInput(input)
	if err != nil {
		return err
	}

	b.DatePublishedEnd = t

	return nil
}

func parseDateInput(input string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:14 -07:00", input)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("2006-01-02 15:14", input)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("2006-01-02", input)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("2006-01", input)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("2006", input)
	if err == nil {
		return t, nil
	}

	return t, err
}
