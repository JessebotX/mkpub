package book

type Series struct {
	Name string

	AlternateNames  []string
	PrimaryNumber   int
	SecondaryNumber int
	About           Content
	Assets          []Asset
	Extra           map[string]any
}
