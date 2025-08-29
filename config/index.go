package config

import (
	"path/filepath"
)

// Index is the main object that contains all [Book]s and
// [SeriesIndex]es.
type Index struct {
	Title            string
	TitleAlternate   []string
	ShortDescription string
	LanguageCode     string
	URL              string
	FaviconImageName string
	Params           map[string]any

	LayoutsDirectory string
	InputPath        string
	Books            []Book
	Series           []SeriesIndex
	Profiles         []Profile
}

func (i *Index) SetDefaults(inputPath string) error {
	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return err
	}
	i.InputPath = absInputPath

	i.Title = filepath.Base(absInputPath)
	i.LayoutsDirectory = filepath.Join(inputPath, "layout")

	return nil
}
