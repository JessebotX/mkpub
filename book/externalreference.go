package book

type ExternalReference struct {
	Address string

	Name               string
	IsAddressHyperlink bool
	Mirrors            []ExternalReference
	Extra              map[string]any
}
