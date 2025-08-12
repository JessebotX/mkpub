package mkpub

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

const (
	IndexConfigName   = "mkpub.yml"
	BookConfigName    = "mkpub-book.yml"
	BookNavConfigName = "mkpub-book-nav.yml"
)

var (
	ErrChapterBookNil                   = errors.New("chapter's book/parent does not exist")
	ErrChapterMissingPossibleIdentifier = errors.New("one of the following values must be defined: \"fileName\", \"title\", \"uniqueID\"")
	ErrContentParsedNil                 = errors.New("parsed content map is not initialized")
	ErrContentFormatNotFound            = errors.New("parsed content format does not exist")
)

func (i *Index) InitDefaults(inputPath string) error {
	i.InputPath = inputPath

	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return err
	}

	i.Title = filepath.Base(absInputPath)
	i.LayoutsDirectory = filepath.Join(inputPath, "layout")

	return nil
}

func (b *Book) InitDefaults(inputPath string, parent *Index) error {
	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return err
	}

	b.InputPath = absInputPath

	b.UniqueID = filepath.Base(b.InputPath)
	b.Title = b.UniqueID
	b.Parent = parent

	if parent != nil {
		b.LanguageCode = parent.LanguageCode

		if parent.URL != "" {
			b.URL, _ = url.JoinPath(parent.URL, "books", b.UniqueID)
		}
	}

	return nil
}

func (b *Book) ChaptersFlattened() []*Chapter {
	var flattened []*Chapter

	for i := range b.Chapters {
		c := &b.Chapters[i]

		var nested []*Chapter
		if len(c.Chapters) > 0 {
			nested = c.ChaptersFlattened()
		}

		flattened = append(flattened, c)
		flattened = append(flattened, nested...)
	}

	return flattened
}

func (c *Chapter) InitDefaults(inputPath string, book *Book) error {
	if book == nil {
		return ErrChapterBookNil
	}

	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return err
	}

	c.InputPath = absInputPath
	c.UniqueID = strings.TrimSuffix(filepath.Base(c.InputPath), ".md")
	c.Title = c.UniqueID
	c.LanguageCode = book.LanguageCode
	c.Book = book

	return nil
}

func (c *Chapter) ChaptersFlattened() []*Chapter {
	var flattened []*Chapter

	for i := range c.Chapters {
		next := &c.Chapters[i]

		var nested []*Chapter
		if len(next.Chapters) > 0 {
			nested = next.ChaptersFlattened()
		}

		flattened = append(flattened, next)
		flattened = append(flattened, nested...)
	}

	return flattened
}

func (s *SeriesIndex) InitDefaults(uniqueID string, parent *Index) {
	s.UniqueID = uniqueID
	s.Parent = parent
}

func (p *Profile) InitDefaults(uniqueID string, parent *Index) {
	p.UniqueID = uniqueID
	p.Parent = parent
}

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

func DecodeIndex(inputPath string) (Index, error) {
	var index Index
	if err := index.InitDefaults(inputPath); err != nil {
		return index, fmt.Errorf("index: failed on initialization: %w", err)
	}

	// --- Unmarshal config file ---
	confBody, err := os.ReadFile(filepath.Join(inputPath, IndexConfigName))
	if err != nil {
		return index, fmt.Errorf("index: failed to read %s: %w", IndexConfigName, err)
	}

	var confMap map[string]any
	if err := yaml.Unmarshal(confBody, &confMap); err != nil {
		return index, fmt.Errorf("index: failed to parse %s: %w", IndexConfigName, err)
	}

	if err := mapToStruct(confMap, &index); err != nil {
		return index, fmt.Errorf("index: failed to parse %s: %w", IndexConfigName, err)
	}

	index.Params = confMap

	// --- Books ---
	booksDir := filepath.Join(inputPath, "books")
	dirs, err := os.ReadDir(booksDir)
	if err != nil {
		return index, fmt.Errorf("index: failed to read books directory: %w", err)
	}

	for i := range dirs {
		dir := dirs[i]

		if !dir.IsDir() {
			continue
		}

		book, err := DecodeBook(filepath.Join(booksDir, dir.Name()), &index)
		if err != nil {
			return index, err
		}

		index.Books = append(index.Books, book)
	}

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

				var output SeriesIndex
				output.InitDefaults(id, &index)
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

				var output Profile
				output.InitDefaults(id, &index)
				output = *author
				output.Content.Raw = []byte(author.About)
				output.Books = append(output.Books, book)

				index.Profiles = append(index.Profiles, output)
			}
		}
	}

	return index, nil
}

func DecodeBook(inputPath string, parent *Index) (Book, error) {
	var book Book
	if err := book.InitDefaults(inputPath, parent); err != nil {
		return book, fmt.Errorf("book \"%s\": failed on initialization: %w", filepath.Base(inputPath), err)
	}

	// --- Unmarshal config file ---
	confBody, err := os.ReadFile(filepath.Join(inputPath, BookConfigName))
	if err != nil {
		return book, fmt.Errorf("book \"%s\": failed to read %s: %w", book.UniqueID, BookConfigName, err)
	}

	var confMap map[string]any
	if err := yaml.Unmarshal(confBody, &confMap); err != nil {
		return book, fmt.Errorf("book \"%s\": failed to parse %s: %w", book.UniqueID, BookConfigName, err)
	}

	if err := mapToStruct(confMap, &book); err != nil {
		return book, fmt.Errorf("book \"%s\": failed to parse %s: %w", book.UniqueID, BookConfigName, err)
	}

	// --- Further parsing ---
	book.Params = confMap
	book.Content.Raw = []byte(book.About)

	if book.Status == "" {
		book.Status = StatusCompleted
	}

	if ok := book.Status.Valid(); !ok {
		return book, fmt.Errorf("book \"%s\": unrecognized status \"%s\". Must be one of the following (case-insensitive): %v", book.UniqueID, book.Status, StatusValidValues)
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

	var chapters []Chapter
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
func parseNav(chapters *[]Chapter, chaptersDir string, book *Book) ([]*Chapter, error) {
	var flattenedList []*Chapter

	for i := range *chapters {
		c := &((*chapters)[i])
		c.InitDefaults(filepath.Join(chaptersDir, c.FileName), book)

		if c.FileName == "" && c.Title == "" && c.UniqueID == "" {
			return nil, ErrChapterMissingPossibleIdentifier
		}

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

		var nested []*Chapter
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
