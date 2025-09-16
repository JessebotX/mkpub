package pub

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

const (
	defaultFilePerms = 0666
	defaultDirPerms  = 0755
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

func WriteBookToStaticSite(book *Book, inputDir, outputDir, layoutsDir string) error {
	if err := os.MkdirAll(outputDir, defaultDirPerms); err != nil {
		return fmt.Errorf("[WRITE BOOK] \"%s\": %w", inputDir, err)
	}

	// --- Parse content ---
	parsedHTML, err := convertMarkdownToHTML(book.Content.Raw)
	if err != nil {
		return writeErrHTMLAndReturn(fmt.Errorf("[WRITE BOOK] \"%s\": %w", inputDir, err), outputDir)
	}
	book.Content.AddFormat("html", parsedHTML)

	// --- Templates ---
	tplName := "index.html"
	tpl, err := template.New("index.html").Funcs(TplFuncs).ParseFiles(filepath.Join(layoutsDir, tplName))
	if err != nil {
		return writeErrHTMLAndReturn(fmt.Errorf("[WRITE BOOK] \"%s\": %w", inputDir, err), outputDir)
	}

	chapterTplName := filepath.Join("_chapter", "index.html")
	chapterTpl, err := template.New("index.html").Funcs(TplFuncs).ParseFiles(filepath.Join(layoutsDir, chapterTplName))
	if err != nil {
		return writeErrHTMLAndReturn(fmt.Errorf("[WRITE BOOK] \"%s\": %w", inputDir, err), outputDir)
	}

	// --- Copy static layout files ---
	if err := copyDirectory(layoutsDir, outputDir, []string{
		tplName,
		chapterTplName,
	}); err != nil {
		return writeErrHTMLAndReturn(fmt.Errorf("[WRITE BOOK] \"%s\": %w", inputDir, err), outputDir)
	}

	// --- Book index.html ---
	f, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		return writeErrHTMLAndReturn(fmt.Errorf("[WRITE BOOK] \"%s\": %w", inputDir, err), outputDir)
	}
	defer f.Close()

	if err := tpl.Execute(f, book); err != nil {
		return writeErrHTMLAndReturn(fmt.Errorf("[WRITE BOOK] \"%s\": %w", inputDir, err), outputDir)
	}

	// --- Chapters ---
	chaptersDir := filepath.Join(outputDir, "chapters")
	if err := os.MkdirAll(chaptersDir, defaultDirPerms); err != nil {
		return fmt.Errorf("[WRITE BOOK] \"%s\": %w", inputDir, err)
	}

	for _, chapter := range book.ChaptersAndSubchapters() {
		if err := writeChapterToStaticSite(chapter, chapter.InputPath, filepath.Join(chaptersDir, chapter.UniqueID+".html"), chapterTpl); err != nil {
			return writeErrHTMLAndReturn(err, outputDir)
		}
	}

	return nil
}

func writeChapterToStaticSite(chapter *Chapter, inputPath, outputPath string, tpl *template.Template) error {
	parsedHTML, err := convertMarkdownToHTML(chapter.Content.Raw)
	if err != nil {
		return fmt.Errorf("[WRITE CHAPTER] \"%s\": %w", inputPath, err)
	}
	chapter.Content.AddFormat("html", parsedHTML)

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("[WRITE CHAPTER] \"%s\": %w", inputPath, err)
	}
	defer f.Close()

	if err := tpl.Execute(f, chapter); err != nil {
		return fmt.Errorf("[WRITE CHAPTER] \"%s\": %w", inputPath, err)
	}

	return nil
}

func convertMarkdownToHTML(rawText []byte) (template.HTML, error) {
	var buffer bytes.Buffer
	if err := md.Convert(rawText, &buffer); err != nil {
		return template.HTML(""), err
	}

	return template.HTML(buffer.String()), nil
}

// Copy files, directories and subdirectories, and supports excluding certain files from copying
func copyDirectory(sourcePath, destinationPath string, excludePaths []string) error {
	return copyDirectoryHelper(sourcePath, destinationPath, sourcePath, excludePaths)
}

// Recursively copies directories and subdirectories. Should use the copyDirectory() function wrapper instead.
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
			if err := copyDirectoryHelper(target, dest, start, excludePaths); err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(newFilePath, defaultDirPerms); err != nil {
				return err
			}

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

// Use like so: (should only be used by functions that are exported so it doesn't write call this function multiple times)
//
//	return fmt.Errorf("<msg>: %w", err))
func writeErrHTMLAndReturn(err error, outputDir string) error {
	_ = filepath.WalkDir(outputDir, func(path string, d os.DirEntry, walkDirErr error) error {
		if walkDirErr != nil {
			return walkDirErr
		}

		if !d.IsDir() && filepath.Ext(path) == ".html" {
			_ = os.WriteFile(path, []byte(`
<div style="
	position:    fixed;
	font-family: ui-monospace, SFMono-Regular, Consolas, 'Liberation Mono', Menlo, monospace;
	top:         0;
	left:        0;
	height:      100%;
	font-weight: bold;
	font-size:   1.75rem;
	box-sizing:  border-box;
	padding:     2rem;
	text-align:  center;
	width:       100%;
	margin:      0;
	padding:     0;
	background:  #333333;
	color:       white;
	z-index:     999">
	<span style="background:red;color:white;">ERROR:</span>`+err.Error()+`</div>`), defaultFilePerms)
		}

		return nil
	})

	return err
}
