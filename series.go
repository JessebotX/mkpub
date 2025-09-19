package pub

import (
	"errors"
	"fmt"
)

var (
	ErrSeriesMissingTitle = errors.New("series: missing title")
)

// Series describes a [Book]'s relation to a set of other [Book] objects (i.e. prequels, sequels, side stories, sharing the same world/universe, etc.)
type Series struct {
	Title           string      `json:"title"`
	TitlesAlternate []string    `json:"titles_alternate"`
	NumberPrimary   uint64      `json:"number_primary"`
	NumberSecondary uint64      `json:"number_secondary"`
	Description     string      `json:"description"`
	Content         Content     `json:"content"`
	External        []Reference `json:"external"`
}

func (s Series) EnsureValid() error {
	if s.Title == "" {
		return ErrSeriesMissingTitle
	}

	for _, e := range s.External {
		if err := e.EnsureValid(); err != nil {
			return fmt.Errorf("series \"%s\": %w", s.Title, err)
		}
	}

	return nil
}

// Compare compares to Series entries based on Series's NumberPrimary and NumberSecondary fields. It returns 0 when they are the same, > 0 when s is higher than other, and < 0 when s is less than other.
func (s Series) Compare(other Series) int {
	if s.NumberPrimary != other.NumberPrimary {
		return int(s.NumberPrimary) - int(other.NumberPrimary)
	}

	return int(s.NumberSecondary) - int(other.NumberSecondary)
}
