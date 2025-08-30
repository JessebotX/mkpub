package mkpub

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/JessebotX/mkpub/config"
)

const (
	IndexConfigName   = "mkpub.yml"
	BookConfigName    = "book.yml"
	BookNavConfigName = "nav.yml"
	IndexBooksDir     = "books"
)

func DecodeIndex(inputPath string) (config.Index, error) {
	var index config.Index

	// --- Unmarshal config file ---
	confMap, err := yamlFileToMap(filepath.Join(inputPath, IndexConfigName))
	if err != nil {
		return index, fmt.Errorf("index: failed to read %s: %w", IndexConfigName, err)
	}

	if err := mapToStruct(confMap, &index); err != nil {
		return index, fmt.Errorf("index: failed to parse %s: %w", IndexConfigName, err)
	}

	if err := index.SetDefaults(inputPath); err != nil {
		return index, fmt.Errorf("index: failed on initialization: %w", err)
	}

	// --- Books ---
	books, err := decodeBooks(filepath.Join(inputPath, IndexBooksDir), &index)
	if err != nil {
		return index, err
	}
	index.Books = books

	// --- Series ---
	for i := range index.Books {
		book := &index.Books[i]

		if len(book.Series) == 0 {
			continue
		}

		for j := range book.Series {
			series := &book.Series[j]

			if series.IndexID == "" && series.Name == "" {
				return index, fmt.Errorf("index: series %d must have either an indexID or a name", j)
			}

			if series.Name == "" {
				series.Name = series.IndexID
			}

			if series.IndexID == "" {
				series.IndexID = series.Name
			}

			exists := false
			for k := range index.Series {
				if series.IndexID == index.Series[k].UniqueID {
					series.SeriesInfo = index.Series[k].SeriesInfo
					index.Series[k].Books = append(index.Series[k].Books, book)

					exists = true
					break
				}
			}

			if !exists {
				id := series.IndexID
				if id == "" {
					id = series.Name
				}

				var output config.SeriesIndex
				output.SetDefaults(id, &index)
				output.SeriesInfo = series.SeriesInfo
				output.Content.Raw = []byte(series.About)
				output.Books = append(output.Books, book)

				index.Series = append(index.Series, output)
			}
		}
	}

	// --- Authors ---
	for i := range index.Books {
		book := &index.Books[i]

		if len(book.Authors) == 0 {
			continue
		}

		for j := range book.Authors {
			author := &book.Authors[j]

			if author.UniqueID == "" && author.Name == "" {
				return index, fmt.Errorf("index: author %d must have either an indexID or a name", j)
			}

			if author.Name == "" {
				author.Name = author.UniqueID
			}

			if author.UniqueID == "" {
				author.UniqueID = author.Name
			}

			exists := false
			for k := range index.Profiles {
				if author.UniqueID == index.Profiles[k].UniqueID {
					*author = index.Profiles[k]
					index.Profiles[k].Books = append(index.Profiles[k].Books, book)

					exists = true
					break
				}
			}

			if !exists {
				id := author.UniqueID
				if id == "" {
					id = author.Name
				}

				var output config.Profile
				output.SetDefaults(id, &index)
				output = *author
				output.Content.Raw = []byte(author.About)
				output.Books = append(output.Books, book)

				index.Profiles = append(index.Profiles, output)
			}
		}
	}

	return index, nil
}

func decodeBooks(booksDir string, index *config.Index) ([]config.Book, error) {
	var books []config.Book

	dirs, err := os.ReadDir(booksDir)
	if err != nil {
		return books, fmt.Errorf("index: failed to read books directory: %w", err)
	}

	for i := range dirs {
		dir := dirs[i]

		if !dir.IsDir() {
			continue
		}

		book, err := DecodeBook(filepath.Join(booksDir, dir.Name()), index)
		if err != nil {
			return books, err
		}

		books = append(books, book)
	}

	return books, nil
}

func DecodeBook(inputPath string, parent *config.Index) (config.Book, error) {
	var book config.Book

	// --- Unmarshal config file ---
	confMap, err := yamlFileToMap(filepath.Join(inputPath, BookConfigName))
	if err != nil {
		return book, fmt.Errorf("book \"%s\": failed to parse %s: %w", book.UniqueID, BookConfigName, err)
	}

	if err := mapToStruct(confMap, &book); err != nil {
		return book, fmt.Errorf("book \"%s\": failed to parse %s: %w", book.UniqueID, BookConfigName, err)
	}

	// --- Further parsing ---
	book.Content.Raw = []byte(book.About)

	if err := book.SetDatePublishedStartFromMap("publishedstart", confMap); err != nil && !errors.Is(err, config.ErrDateFromMapKeyNotFound) {
		return book, fmt.Errorf("book \"%s\": %w", book.UniqueID, err)
	}

	if err := book.SetDatePublishedEndFromMap("publishedend", confMap); err != nil && !errors.Is(err, config.ErrDateFromMapKeyNotFound) {
		return book, fmt.Errorf("book \"%s\": %w", book.UniqueID, err)
	}

	if err := book.SetDefaults(inputPath, parent); err != nil {
		return book, fmt.Errorf("book \"%s\": failed on initialization: %w", filepath.Base(inputPath), err)
	}

	// --- Parse chapters ---
	navBody, err := os.ReadFile(filepath.Join(inputPath, BookNavConfigName))
	if err != nil {
		return book, fmt.Errorf("book \"%s\": failed to read %s: %w", book.UniqueID, BookNavConfigName, err)
	}

	var navConfMap []any
	if err := yaml.Unmarshal(navBody, &navConfMap); err != nil {
		return book, fmt.Errorf("book \"%s\": failed to parse %s: %w", book.UniqueID, BookNavConfigName, err)
	}

	var chapters []config.Chapter
	if err := mapToStruct(navConfMap, &chapters); err != nil {
		return book, fmt.Errorf("book \"%s\": failed to parse %s: %w", book.UniqueID, BookNavConfigName, err)
	}

	chaptersDir := filepath.Join(inputPath, "chapters")
	flattenedChapters, err := parseNav(&chapters, chaptersDir, &book)
	if err != nil {
		return book, fmt.Errorf("book \"%s\": %w", book.UniqueID, err)
	}

	// Set next and previous values
	for i := range flattenedChapters {
		c := flattenedChapters[i]

		if i-1 >= 0 {
			c.Previous = flattenedChapters[i-1]
		}

		if i+1 < len(flattenedChapters) {
			c.Next = flattenedChapters[i+1]
		}
	}

	book.Chapters = chapters

	return book, nil
}

// Returns a list of chapters in a flattened array for the purposes of pagination order.
func parseNav(chapters *[]config.Chapter, chaptersDir string, book *config.Book) ([]*config.Chapter, error) {
	var flattenedList []*config.Chapter

	for i := range *chapters {
		c := &((*chapters)[i])

		inputPath := filepath.Join(chaptersDir, c.FileName)
		absInputPath, err := filepath.Abs(inputPath)
		if err != nil {
			return nil, err
		}

		c.InputPath = absInputPath
		c.UniqueID = strings.TrimSuffix(filepath.Base(c.InputPath), ".md")
		if c.LanguageCode == "" {
			c.LanguageCode = book.LanguageCode
		}
		c.Book = book

		// if c.FileName == "" && c.Title == "" && c.UniqueID == "" {
		// 	return nil, ErrChapterMissingPossibleIdentifier
		// }

		// If either title or uniqueID is missing...
		if c.Title != "" && c.UniqueID == "" {
			c.UniqueID = c.Title
		} else if c.Title == "" && c.UniqueID != "" {
			c.Title = c.UniqueID
		}

		if c.FileName != "" {
			raw, err := os.ReadFile(filepath.Join(chaptersDir, c.FileName))
			if err != nil {
				return nil, err
			}
			c.Content.Raw = raw
			c.InputPath = filepath.Join(chaptersDir, c.FileName)

			// if title or uniqueID is still missing, use fileName without the file extension
			if c.UniqueID == "" {
				c.UniqueID = strings.TrimSuffix(c.FileName, ".md")
			}

			if c.Title == "" {
				c.Title = strings.TrimSuffix(c.FileName, ".md")
			}
		}

		var nested []*config.Chapter
		if c.Chapters != nil {
			l, err := parseNav(&c.Chapters, chaptersDir, book)
			if err != nil {
				return nil, err
			}
			nested = l
		}

		flattenedList = append(flattenedList, c)
		flattenedList = append(flattenedList, nested...)
	}

	return flattenedList, nil
}

func yamlFileToMap(configPath string) (map[string]any, error) {
	var conf map[string]any

	f, err := os.ReadFile(configPath)
	if err != nil {
		return conf, err
	}

	if err := yaml.Unmarshal(f, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}

func mapToStruct[M map[string]any | []any](m M, s any) error {
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, s); err != nil {
		return err
	}

	return nil
}
