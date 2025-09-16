package pub

import (
	"errors"
	"fmt"
)

var (
	ErrSeriesMissingTitles = errors.New("series: missing at least 1 title in Titles")
)

type ErrSeriesEmptyTitle struct {
	EntryNumber int // starts at 1 instead of 0
}

func (e ErrSeriesEmptyTitle) Error() string {
	return fmt.Sprintf("series: invalid empty title in Titles entry number %v", e.EntryNumber)
}

// Series describes a [Book]'s relation to a set of other [Book]s (i.e. prequels, sequels, side stories, sharing the same world/universe, etc.)
type Series struct {
	Titles          []string    `json:"titles"`
	NumberPrimary   uint64      `json:"number_primary"`
	NumberSecondary uint64      `json:"number_secondary"`
	Description     string      `json:"description"`
	Content         Content     `json:"content"`
	External        []Reference `json:"external"`
}

func (s Series) Title() string {
	return s.Titles[0]
}

func (s Series) EnsureValid() error {
	if len(s.Titles) == 0 {
		return ErrSeriesMissingTitles
	}

	for i, title := range s.Titles {
		if title == "" {
			return ErrSeriesEmptyTitle{EntryNumber: i + 1}
		}
	}

	for _, e := range s.External {
		if err := e.EnsureValid(); err != nil {
			return fmt.Errorf("series \"%s\": %w", s.Title(), err)
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
