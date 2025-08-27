package book

type Content struct {
	Raw []byte

	parsed map[string]any
}

func (c *Content) New(raw []byte) {
	c.Raw = raw
	if c.parsed == nil {
		c.parsed = make(map[string]any)
	}
}

func (c *Content) Format(key string) (any, bool) {
	result, ok := c.parsed[key]
	if !ok {
		return "", false
	}

	return result, true
}

func (c *Content) AddFormat(key string, content any) {
	if c.parsed == nil {
		c.parsed = make(map[string]any, 1)
	}

	c.parsed[key] = content
}
