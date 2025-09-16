package pub

import (
	"errors"
	"fmt"
)

var (
	ErrReferenceMissingName = errors.New("reference: missing Name")
)

type ErrReferenceMissingAddress struct {
	Name string
}

func (e ErrReferenceMissingAddress) Error() string {
	return fmt.Sprintf("reference: missing address for reference with name \"%s\"", e.Name)
}

type ErrReferenceTooManyNestedDomains struct {
	TopLevelName        string
	DomainAlternateName string
}

func (e ErrReferenceTooManyNestedDomains) Error() string {
	return fmt.Sprintf("reference\"%s\": DomainAlternate \"%s\" cannot have its own DomainsAlternate references (too many nested levels)", e.TopLevelName, e.DomainAlternateName)
}

// Reference represents an external link/address that is generally clickable.
type Reference struct {
	Name             string      `json:"name"`
	Address          string      `json:"address"`
	NotClickable     bool        `json:"not_clickable"`
	DomainsAlternate []Reference `json:"domains_alternate"`
}

func (r *Reference) EnsureValid() error {
	return r.checkValid(r.Name, 0)
}

// Level starts at 0. Can only nest up to 2 levels (level either 0 or 1 or it returns err)
func (r *Reference) checkValid(topLevelName string, level int) error {
	if r.Name == "" && r.Address != "" {
		r.Name = r.Address
	}

	if r.Name == "" {
		return ErrReferenceMissingName
	}

	if r.Address == "" {
		return ErrReferenceMissingAddress{Name: r.Name}
	}

	if level > 1 {
		return ErrReferenceTooManyNestedDomains{TopLevelName: topLevelName, DomainAlternateName: r.Name}
	}
	for _, alt := range r.DomainsAlternate {
		if err := alt.checkValid(topLevelName, level+1); err != nil {
			return err
		}
	}

	return nil
}
