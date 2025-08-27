package book

type Asset struct {
	Name string

	Type          string
	Subtype       string
	AlternateText string
	Caption       Content
	Extra         map[string]any
}
