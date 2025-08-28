package config

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

// SeriesIndex describes a series: a set of related books.
type SeriesIndex struct {
	SeriesInfo

	UniqueID string
	Parent   *Index
	Books    []*Book
	Content  Content
}

func (s *SeriesIndex) EnsureDefaultsSet(uniqueID string, parent *Index) {
	s.UniqueID = uniqueID
	s.Parent = parent
}
