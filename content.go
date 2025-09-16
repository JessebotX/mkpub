package pub

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	ErrContentParsedUninitialized = errors.New("parsed content map is not initialized")
	ErrContentFormatNotFound      = errors.New("parsed content format does not exist")
)

// Content represents a body of text that is/can be parsed into different formats (e.g. Markdown to HTML, etc.).
type Content struct {
	Raw []byte

	parsed map[string]any
}

func (c *Content) MarshalText() ([]byte, error) {
	return c.Raw, nil
}

func (c *Content) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", string(c.Raw))), nil
}

func (c *Content) MarshalYAML() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", string(c.Raw))), nil
}

func (c *Content) UnmarshalText(text []byte) error {
	c.Raw = text
	return nil
}

func (c *Content) UnmarshalJSON(text []byte) error {
	return c.UnmarshalText(bytes.TrimPrefix(bytes.TrimSuffix(text, []byte("\"")), []byte("\"")))
}

func (c *Content) UnmarshalYAML(text []byte) error {
	return c.UnmarshalText(bytes.TrimPrefix(bytes.TrimSuffix(text, []byte("\"")), []byte("\"")))
}

func (c *Content) Format(format string) (any, error) {
	if c.parsed == nil {
		return "", ErrContentParsedUninitialized
	}

	res, ok := c.parsed[format]
	if !ok {
		return "", ErrContentFormatNotFound
	}

	return res, nil
}

func (c *Content) AddFormat(format string, content any) {
	if c.parsed == nil {
		c.parsed = make(map[string]any, 1)
	}

	c.parsed[format] = content
}
