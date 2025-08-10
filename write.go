package mkpub

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
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
		return fmt.Errorf("write index: failed to copy files to output: %w", err)
	}

	// --- book index ---
	wrIndex, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		return fmt.Errorf("write index: %w", err)
	}
	defer wrIndex.Close()

	indexTmplPath := filepath.Join(index.LayoutsDirectory, "index.html")
	indexTmpl, err := template.New("index.html").Funcs(TemplateFuncs).ParseFiles(indexTmplPath)
	if err != nil {
		err = fmt.Errorf("write index: failed to read index template file %s: %w", indexTmplPath, err)
		writeErrHTML(wrIndex, err)
		return err
	}

	if err := indexTmpl.ExecuteTemplate(wrIndex, "index.html", index); err != nil {
		err = fmt.Errorf("write index: %w", err)
		writeErrHTML(wrIndex, err)
		return err
	}

	for i := range index.Books {
		book := &index.Books[i]
		bookOutputDir := filepath.Join(outputDir, "books", book.UniqueID)

		if err := writeBookToStaticWebsite(book, bookOutputDir); err != nil {
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
		return fmt.Errorf("write \"%s\": %w", book.UniqueID, err)
	}
	defer wrBook.Close()

	index := book.Parent
	bookTmplPath := filepath.Join(index.LayoutsDirectory, "_book.html")
	bookTmpl, err := template.New("_book.html").Funcs(TemplateFuncs).ParseFiles(bookTmplPath)
	if err != nil {
		err = fmt.Errorf("write \"%s\": failed to read template file %s: %w", book.UniqueID, bookTmplPath, err)
		writeErrHTML(wrBook, err)
		return err
	}

	if err := bookTmpl.ExecuteTemplate(wrBook, "_book.html", book); err != nil {
		err = fmt.Errorf("write \"%s\": %w", book.UniqueID, err)
		writeErrHTML(wrBook, err)
		return err
	}

	for i := range book.Chapters {
		chapter := &book.Chapters[i]
		chapterOutputDir := filepath.Join(outputDir, "chapters")
		if err := os.MkdirAll(chapterOutputDir, 0755); err != nil {
			return err
		}

		if err := writeChapterToStaticWebsite(chapter, filepath.Join(chapterOutputDir, chapter.UniqueID+".html"), book); err != nil {
			return err
		}
	}

	return nil
}

func writeChapterToStaticWebsite(chapter *OutputChapter, outputPath string, book *OutputBook) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	index := book.Parent
	chapterTmplPath := filepath.Join(index.LayoutsDirectory, "_chapter.html")
	chapterTmpl, err := template.New("_chapter.html").Funcs(TemplateFuncs).ParseFiles(chapterTmplPath)
	if err != nil {
		err = fmt.Errorf("write \"%s\" chapter \"%s\": failed to read template file %s: %w", book.UniqueID, chapter.UniqueID, chapterTmplPath, err)
		writeErrHTML(f, err)
		return err
	}

	if err := chapterTmpl.ExecuteTemplate(f, "_chapter.html", chapter); err != nil {
		err = fmt.Errorf("write \"%s\" chapter \"%s\": %w", book.UniqueID, chapter.UniqueID, err)
		writeErrHTML(f, err)
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

func writeErrHTML(wr io.Writer, err error) {
	_, _ = wr.Write([]byte(`<div style="position:fixed;
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
