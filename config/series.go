package config

import (
	"errors"
	"strings"
)

var ErrSeriesMissingPossibleIdentifier = errors.New("book series must either have an indexID or a name")

// SeriesInfo describes internal information of a series
type SeriesInfo struct {
	Name             string
	NameAlternate    []string
	ShortDescription string
	About            string
	IDs              map[string]string
	Links            []ExternalReference
	Images           []Asset
}

// BookSeriesItem describes a [Book]'s relation/entry in a series.
type BookSeriesItem struct {
	SeriesInfo

	IndexID     string
	EntryNumber float64
}

func (b *BookSeriesItem) SetDefaults() {
	if b.Name == "" {
		b.Name = b.IndexID
	}

	if b.IndexID == "" {
		b.IndexID = b.Name
	}
}

func (b *BookSeriesItem) IsValid() error {
	if b.IndexID == "" && b.Name == "" {
		return ErrSeriesMissingPossibleIdentifier
	}

	return nil
}

// SeriesIndex describes a series: a set of related books.
type SeriesIndex struct {
	SeriesInfo

	UniqueID string
	Parent   *Index
	Books    []*Book
	Content  Content
}

func (s *SeriesIndex) SetDefaults(uniqueID string, parent *Index) {
	if s.Parent == nil {
		s.Parent = parent
	}

	uniqueID = strings.TrimSpace(uniqueID)
	if s.UniqueID == "" {
		s.UniqueID = uniqueID
	}
}
