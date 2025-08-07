package mkpub

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
		StatusInactive,
		StatusInactive,
		StatusOngoing,
	}
)

// ExternalReference points to an external object, such as a hyperlink
// that directs the user to external donation pages.
type ExternalReference struct {
	Name          string
	Address       string
	IsHyperlink   bool
	IconImagePath string
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
	ImagePaths       []string
	Links            []ExternalReference
}

// SeriesItem describes a [Book]'s relation/entry in a series.
type SeriesItem struct {
	Name             string
	NameAlternate    []string
	ShortDescription string
	About            string
	IDs              map[string]string
	Links            []ExternalReference
	ImagePaths       []string

	SeriesIndexID string
	EntryNumber   float64
}

// SeriesIndex describes a series: a set of related books.
type SeriesIndex struct {
	Name             string
	NameAlternate    []string
	ShortDescription string
	IDs              map[string]string
	Links            []ExternalReference
	ImagePaths       []string

	UniqueID string
}

// Index is the main object that contains all [Book]s and
// [SeriesIndex]es.
type Index struct {
	Title            string
	TitleAlternate   []string
	ShortDescription string
	LanguageCode     string
	URL              string
	Profiles         []Profile
	Series           []SeriesIndex
	FaviconImagePath string
	Params           map[string]any

	LayoutsDirectory string
}

// Book is a written work that contains one or more [Chapter]s, which are
// usually defined in a specific reading order. It may also be apart of zero or
// more [SeriesItem]'s.
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
	CoverImagePath     string
	ImagePaths         []string
	Copyright          string
	License            string
	DatePublishedStart string
	DatePublishedEnd   string
	IDs                map[string]string
	Params             map[string]any

	CharacterEncoding string
	UniqueID          string
}

// Chapter represents a division in a [Book].
type Chapter struct {
	Title            string
	TitleAlternate   []string
	Subtitle         string
	TitleSort        string
	ShortDescription string
	Content          string
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

	UniqueID string
}
