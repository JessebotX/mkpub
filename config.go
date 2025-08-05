package mkpub

const (
	StatusCompleted = "completed"
	StatusInactive  = "inactive"
	StatusOngoing   = "ongoing"
	StatusHiatus    = "hiatus"
)

var (
	StatusValidValues = []string{
		StatusCompleted,
		StatusInactive,
		StatusInactive,
		StatusOngoing,
	}
)

type ExternalReference struct {
	Name          string
	Address       string
	IsHyperlink   bool
	IconImagePath string
}

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

type SeriesIndex struct {
	Name             string
	NameAlternate    []string
	ShortDescription string
	IDs              map[string]string
	Links            []ExternalReference
	ImagePaths       []string

	UniqueID string
}

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

	InputPath        string
	LayoutsDirectory string
}

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
	Status             string
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

type Chapter struct {
	Title            string
	TitleAlternate   []string
	Subtitle         string
	TitleSort        string
	ShortDescription string
	Content          string
	AuthorsNote      string
	Status           string
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
