package config

// Profile can represent an individual author or contributor, or an
// organization such as a publisher.
type Profile struct {
	UniqueID         string
	Name             string
	NameAlternate    []string
	Roles            []string
	ShortDescription string
	About            string
	Images           []Asset
	Links            []ExternalReference

	Parent  *Index
	Books   []*Book
	Content Content
}

func (p *Profile) SetDefaults(uniqueID string, parent *Index) {
	p.UniqueID = uniqueID
	p.Parent = parent
}
