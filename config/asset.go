package config

// Asset refers to assets such as images and videos that are included within a [Book].
type Asset struct {
	Name          string
	AlternateText string
	Caption       string
	Type          string
	Format        string
}
