package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/JessebotX/pub"
)

const (
	defaultDirPerms = 0755
)

type BookCommand struct {
	Init BookInitCommand `cmd:"" help:"Initialize a new book project"`
}

type BookInitCommand struct {
	Path string `arg:"" optional:"" default:"./" help:"Path to initialize a book project. Directory path will be automatically created if it does not already exist"`
}

func (b BookInitCommand) Run(ctx *Context) error {
	absPath, err := filepath.Abs(b.Path)
	if err != nil {
		return err
	}

	// Create directories
	if err := os.MkdirAll(absPath, defaultDirPerms); err != nil {
		return err
	}

	chaptersDirPath := filepath.Join(b.Path, pub.BookChaptersDirName)
	if err := os.MkdirAll(chaptersDirPath, defaultDirPerms); err != nil {
		return err
	}

	layoutsDirPath := filepath.Join(b.Path, LayoutsDirName)
	if err := os.MkdirAll(layoutsDirPath, defaultDirPerms); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(layoutsDirPath, "_chapter"), defaultDirPerms); err != nil {
		return err
	}

	assetsDirPath := filepath.Join(b.Path, pub.BookAssetsDirName)
	if err := os.MkdirAll(assetsDirPath, defaultDirPerms); err != nil {
		return err
	}

	// Create layout templates (TODO: add default templates)
	fLayoutIndex, err := os.Create(filepath.Join(layoutsDirPath, "index.html"))
	if err != nil {
		return err
	}
	defer fLayoutIndex.Close()

	if _, err := fLayoutIndex.WriteString(`<!DOCTYPE html>
<html lang="{{ .LanguageCode }}">
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{ .Title }}</title>

<h1>{{ .Title }}</h1>
<p>Nothing here yet...</p>`); err != nil {
		return err
	}

	fLayoutChapter, err := os.Create(filepath.Join(layoutsDirPath, "_chapter", "index.html"))
	if err != nil {
		return err
	}
	defer fLayoutChapter.Close()

	if _, err := fLayoutChapter.WriteString(`<!DOCTYPE html>
<html lang="{{ with .LanguageCode }}{{ . }}{{ else }}{{ .Book.LanguageCode }}{{ end }}">
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{ .Title }} | {{ .Book.Title }}</title>

<h1>{{ .Title }}</h1>
<p>Nothing here yet...</p>`); err != nil {
		return err
	}

	// Create config files
	fPub, err := os.Create(filepath.Join(b.Path, pub.BookConfigFileName))
	if err != nil {
		return err
	}
	defer fPub.Close()

	// TODO: probably interactively prompt user on certain required fields
	if _, err := fPub.Write([]byte(`unique_id: "` + strings.ToLower(filepath.Base(absPath)) + `"`)); err != nil {
		return err
	}

	fChapters, err := os.Create(filepath.Join(b.Path, pub.BookChaptersConfigFileName))
	if err != nil {
		return err
	}
	defer fChapters.Close()

	return nil
}
