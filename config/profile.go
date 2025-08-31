package config

import (
	"errors"
	"strings"
)

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

func (p *Profile) EnsureDefaults() {
	if p.Name == "" {
		p.Name = p.UniqueID
	}

	if p.UniqueID == "" {
		p.UniqueID = strings.TrimSpace(p.Name)
	}
}

var ErrProfileMissingIdentifier = errors.New("author must either have an indexID or a name")

func (p *Profile) IsValid() error {
	if p.UniqueID == "" && p.Name == "" {
		return ErrProfileMissingIdentifier
	}

	return nil
}
