package pub

import (
	"errors"
	"fmt"
)

var (
	ErrProfileMissingNames = errors.New("profile: missing 1 or more Names")
)

type ErrProfileEmptyName struct {
	Index int // starts at 1 instead of 0
}

func (e ErrProfileEmptyName) Error() string {
	return fmt.Sprintf("profile: invalid empty name in Names entry number %d", e.Index)
}

// Profile may represent an individual or an organization that is credited as either an author, contributor or publisher affliated with a [Book].
type Profile struct {
	Names    []string    `json:"names"`
	NameSort string      `json:"name_sort"`
	Content  Content     `json:"content"`
	External []Reference `json:"external"`
}

func (p *Profile) Name() string {
	return p.Names[0]
}

func (p *Profile) EnsureValid() error {
	if len(p.Names) == 0 {
		return ErrProfileMissingNames
	}

	for i, name := range p.Names {
		if name == "" {
			return ErrProfileEmptyName{Index: i + 1}
		}
	}

	if p.NameSort == "" {
		p.NameSort = p.Name()
	}

	for _, e := range p.External {
		if err := e.EnsureValid(); err != nil {
			return fmt.Errorf("profile \"%s\": %w", p.Name(), err)
		}
	}

	return nil
}
