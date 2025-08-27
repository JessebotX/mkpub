package book

type Profile struct {
	Name string

	AlternateNames []string
	About          Content
	Links          []ExternalReference
	Image          Asset
	OtherAssets    []Asset
	Extra          map[string]any
}
