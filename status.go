package pub

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

type Status int

const (
	Completed Status = iota
	Ongoing
	Hiatus
	Inactive
)

var (
	StatusMap = map[string]Status{
		"completed": Completed,
		"ongoing":   Ongoing,
		"hiatus":    Hiatus,
		"inactive":  Inactive,
	}
)

type ErrStatusMarshalUnrecognized struct {
	Value Status
}

func (e ErrStatusMarshalUnrecognized) Error() string {
	return fmt.Sprintf("status: unrecognized value %d", int(e.Value))
}

type ErrStatusUnmarshalUnrecognized struct {
	StatusString string
}

func (e ErrStatusUnmarshalUnrecognized) Error() string {
	return fmt.Sprintf("status: unrecognized value \"%s\" (value must be one of the following (case doesn't matter): %v)", e.StatusString, strings.Join(slices.Sorted(maps.Keys(StatusMap)), ", "))
}

func (s *Status) UnmarshalText(text []byte) error {
	vStr := string(text)

	v, ok := StatusMap[strings.ToLower(vStr)]
	if !ok {
		return ErrStatusUnmarshalUnrecognized{StatusString: vStr}
	}

	*s = v

	return nil
}

func (s Status) MarshalText() ([]byte, error) {
	for k, v := range StatusMap {
		if v == s {
			return []byte(k), nil
		}
	}

	return nil, ErrStatusMarshalUnrecognized{Value: s}
}

func (s Status) String() (string, error) {
	b, err := s.MarshalText()
	if err != nil {
		return "", err
	}

	return string(b), nil
}
