package config

import "errors"

var (
	ErrContentParsedNil      = errors.New("parsed content map is not initialized")
	ErrContentFormatNotFound = errors.New("parsed content format does not exist")
)

type Content struct {
	Raw []byte

	parsed map[string]any
}

func (c *Content) Format(format string) (any, error) {
	if c.parsed == nil {
		return "", ErrContentParsedNil
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
