package config

// ExternalReference points to an external object, such as a hyperlink
// that directs the user to external donation pages.
type ExternalReference struct {
	Name        string
	Address     string
	IsHyperlink bool

	Mirrors []ExternalReference
}
