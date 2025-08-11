package mkpub

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	mdhtml "github.com/yuin/goldmark/renderer/html"
)

var (
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Typographer,
		),
		goldmark.WithParserOptions(
			parser.WithAttribute(),
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			mdhtml.WithXHTML(),
		),
	)
)

func WriteIndexToStaticWebsite(index *OutputIndex, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// --- handle static files ---
	if err := copyDirectory(index.LayoutsDirectory, outputDir, []string{
		"index.html",
		"_book.html",
		"_chapter.html",
		"_profile.html",
		"_series.html",
		"_tag.html",
	}); err != nil {
		return fmt.Errorf("write: failed to copy files to output: %w", err)
	}

	// --- parse content ---
	for i := range index.Books {
		book := &index.Books[i]

		parsedHTML, err := convertMarkdownToHTML(book.Content.Raw)
		if err != nil {
			return fmt.Errorf("write: failed to convert book \"%s\" about field to html: %w", book.UniqueID, err)
		}
		book.Content.AddFormat("html", parsedHTML)

		for j := range book.ChaptersFlattened() {
			chapter := book.ChaptersFlattened()[j]

			parsedHTML, err := convertMarkdownToHTML(chapter.Content.Raw)
			if err != nil {
				return fmt.Errorf("write: failed to convert chapter \"%s\" content to html: %w", chapter.UniqueID, err)
			}
			chapter.Content.AddFormat("html", parsedHTML)
		}
	}

	for i := range index.Series {
		series := &index.Series[i]
		parsedHTML, err := convertMarkdownToHTML(series.Content.Raw)
		if err != nil {
			return fmt.Errorf("write: failed to convert series \"%s\" (%s) about info to html: %w", series.Name, series.UniqueID, err)
		}
		series.Content.AddFormat("html", parsedHTML)
	}

	for i := range index.Profiles {
		profile := &index.Profiles[i]
		parsedHTML, err := convertMarkdownToHTML(profile.Content.Raw)
		if err != nil {
			return fmt.Errorf("write: failed to convert profile \"%s\" (%s) about info to html: %w", profile.Name, profile.UniqueID, err)
		}
		profile.Content.AddFormat("html", parsedHTML)
	}

	// --- favicon ---
	faviconName := index.FaviconImageName
	if faviconName != "" {
		if err := copyFile(filepath.Join(index.InputPath, faviconName), filepath.Join(outputDir, faviconName)); err != nil {
			return err
		}
	}

	// --- book index ---
	wrIndex, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	defer wrIndex.Close()

	indexTmplPath := filepath.Join(index.LayoutsDirectory, "index.html")
	indexTmpl, err := template.New("index.html").Funcs(TemplateFuncs).ParseFiles(indexTmplPath)
	if err != nil {
		err = fmt.Errorf("write: failed to read index template file %s: %w", indexTmplPath, err)
		writeErrHTML(err, wrIndex)
		return err
	}

	if err := indexTmpl.ExecuteTemplate(wrIndex, "index.html", index); err != nil {
		err = fmt.Errorf("write: %w", err)
		writeErrHTML(err, wrIndex)
		return err
	}

	for i := range index.Books {
		book := &index.Books[i]
		bookOutputDir := filepath.Join(outputDir, "books", book.UniqueID)

		if err := writeBookToStaticWebsite(book, bookOutputDir); err != nil {
			writeErrHTML(err, wrIndex)
			return err
		}
	}

	return nil
}

func writeBookToStaticWebsite(book *OutputBook, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	wrBook, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		return fmt.Errorf("book \"%s\" (%s): %w", book.Title, book.UniqueID, err)
	}
	defer wrBook.Close()

	// --- Book main page ---

	index := book.Parent
	bookTmplPath := filepath.Join(index.LayoutsDirectory, "_book.html")
	bookTmpl, err := template.New("_book.html").Funcs(TemplateFuncs).ParseFiles(bookTmplPath)
	if err != nil {
		err = fmt.Errorf("book \"%s\" (%s): failed to read template file %s: %w", book.Title, book.UniqueID, bookTmplPath, err)
		writeErrHTML(err, wrBook)
		return err
	}

	if err := bookTmpl.ExecuteTemplate(wrBook, "_book.html", book); err != nil {
		err = fmt.Errorf("book \"%s\" (%s): %w", book.Title, book.UniqueID, err)
		writeErrHTML(err, wrBook)
		return err
	}

	// --- Cover image ---
	imagesOutputDir := filepath.Join(outputDir, "images")
	if err := os.MkdirAll(imagesOutputDir, 0755); err != nil {
		err = fmt.Errorf("book \"%s\" (%s): %w", book.Title, book.UniqueID, err)
		writeErrHTML(err, wrBook)
		return err
	}

	coverName := book.CoverImage.Name
	if coverName != "" {
		if err := copyFile(filepath.Join(book.InputPath, "images", coverName), filepath.Join(imagesOutputDir, coverName)); err != nil {
			err = fmt.Errorf("book \"%s\" (%s): %w", book.Title, book.UniqueID, err)
			writeErrHTML(err, wrBook)
			return err
		}
	}

	// --- Chapters ---

	flattenedChapters := book.ChaptersFlattened()
	for i := range flattenedChapters {
		chapter := flattenedChapters[i]
		chapterOutputDir := filepath.Join(outputDir, "chapters")
		if err := os.MkdirAll(chapterOutputDir, 0755); err != nil {
			err = fmt.Errorf("book \"%s\" (%s): %w", book.Title, book.UniqueID, err)
			writeErrHTML(err, wrBook)
			return err
		}

		if err := writeChapterToStaticWebsite(chapter, filepath.Join(chapterOutputDir, chapter.UniqueID+".html")); err != nil {
			err = fmt.Errorf("book \"%s\" (%s): %w", book.Title, book.UniqueID, err)
			writeErrHTML(err, wrBook)
			return err
		}
	}

	return nil
}

func writeChapterToStaticWebsite(chapter *OutputChapter, outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	book := chapter.Book
	layoutsDir := book.Parent.LayoutsDirectory
	chapterTmplPath := filepath.Join(layoutsDir, "_chapter.html")
	chapterTmpl, err := template.New("_chapter.html").Funcs(TemplateFuncs).ParseFiles(chapterTmplPath)
	if err != nil {
		err = fmt.Errorf("chapter \"%s\": failed to read template file %s: %w", chapter.UniqueID, chapterTmplPath, err)
		writeErrHTML(err, f)
		return err
	}

	if err := chapterTmpl.ExecuteTemplate(f, "_chapter.html", chapter); err != nil {
		err = fmt.Errorf("chapter \"%s\": %w", chapter.UniqueID, err)
		writeErrHTML(err, f)
		return err
	}

	return nil
}

func copyDirectory(sourcePath, destinationPath string, excludePaths []string) error {
	return copyDirectoryHelper(sourcePath, destinationPath, sourcePath, excludePaths)
}

func copyDirectoryHelper(curr, dest, start string, excludePaths []string) error {
	items, err := os.ReadDir(curr)
	if err != nil {
		return err
	}

	for _, item := range items {
		target := filepath.Join(curr, item.Name())
		targetFromStart := strings.TrimLeft(strings.TrimPrefix(target, start), "/\\")

		// check exclusions
		if slices.Contains(excludePaths, targetFromStart) {
			continue
		}

		newFilePath := filepath.Join(dest, targetFromStart)

		// recursively copy subdirectories
		if item.IsDir() {
			if err := os.MkdirAll(newFilePath, 0755); err != nil {
				return err
			}

			if err := copyDirectoryHelper(target, dest, start, excludePaths); err != nil {
				return err
			}
		} else {
			if err := copyFile(target, newFilePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(sourcePath, destinationPath string) error {
	in, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return nil
}

func writeErrHTML(err error, writers ...io.Writer) {
	for i := range writers {
		_, _ = writers[i].Write([]byte(`<div style="position:fixed;
font-family:ui-monospace,SFMono-Regular,Consolas,'Liberation Mono',Menlo,monospace;
top:0;
left:0;
height:100%;
font-weight:bold;
font-size:1.75rem;
box-sizing:border-box;
padding:2rem;
text-align:center;
width:100%;
margin:0;
padding:0;
background:#333333;
color:white;
z-index:999"><span style="background:red;color:white;">ERROR:</span> ` + err.Error() + `</div>`))
	}
}

func convertMarkdownToHTML(data []byte) (template.HTML, error) {
	var buffer bytes.Buffer
	if err := md.Convert(data, &buffer); err != nil {
		return template.HTML(""), err
	}

	return template.HTML(buffer.String()), nil
}
