package pub

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

var (
	DateTimeLayouts = []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02T15:04Z07:00",
		"2006-01-02 15:04Z07:00",
		"2006-01-02TZ07:00",
		"2006-01-02 Z07:00",
		"2006-01TZ07:00",
		"2006-01 Z07:00",
		"2006TZ07:00",
		"2006 Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006-01",
		"2006",
	}
)

type DateTime struct {
	time.Time
}

func (d DateTime) MarshalText() ([]byte, error) {
	timeString := d.Format("2006-01-02 15:04:05Z07:00")
	return []byte(timeString), nil
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	timeString := fmt.Sprintf("\"%s\"", d.Format("2006-01-02T15:04:05Z07:00"))
	return []byte(timeString), nil
}

func (d DateTime) MarshalYAML() ([]byte, error) {
	timeString := fmt.Sprintf("\"%s\"", d.Format("2006-01-02T15:04:05Z07:00"))
	return []byte(timeString), nil
}

func (d *DateTime) UnmarshalText(text []byte) error {
	input := string(text)

	t, err := dateFromString(input)
	if err != nil {
		return err
	}
	*d = t

	return nil
}

func (d *DateTime) UnmarshalJSON(text []byte) error {
	return d.UnmarshalText(bytes.TrimPrefix(bytes.TrimSuffix(text, []byte("\"")), []byte("\"")))
}

func (d *DateTime) UnmarshalYAML(text []byte) error {
	return d.UnmarshalText(bytes.TrimPrefix(bytes.TrimSuffix(text, []byte("\"")), []byte("\"")))
}

func dateFromString(input string) (DateTime, error) {
	var errs error

	for _, layout := range DateTimeLayouts {
		t, err := time.Parse(layout, input)
		if err == nil {
			return DateTime{t}, nil
		}
		errs = errors.Join(errs, err)
	}

	return DateTime{time.Time{}}, errs
}
