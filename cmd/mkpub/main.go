package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/alecthomas/kong"

	"github.com/JessebotX/mkpub"
)

const (
	Version = "0.12.2"
)

var (
	cli CLI
)

type Context struct {
	NoNonEssentialMessages bool
	Plain                  bool
}

type BuildCommand struct {
	InputDirectory  string `short:"i" help:"Path to directory containing source files." type:"path" default:"./"`
	OutputDirectory string `short:"o" help:"Path to directory containing compiled output files/formats for distribution." type:"path"`
	Minify          bool   `help:"Str output of supported file formats."`
}

func (b *BuildCommand) Run(context *Context) error {
	if b.OutputDirectory == "" {
		b.OutputDirectory = filepath.Join(b.InputDirectory, "build")
	}

	// --- Begin build ---
	buildTimeStart := time.Now()

	if !context.NoNonEssentialMessages {
		minificationStatus := strconv.FormatBool(b.Minify)
		if !b.Minify {
			minificationStatus += " (\"--minify\" to enable) TODO: minification not implemented yet"
		}

		fmt.Println(terminalStyle("Building...", TerminalTextBold, TerminalTextGreen))
		fmt.Printf("+ Input directory:     %s\n", b.InputDirectory)
		fmt.Printf("+ Output directory:    %s\n", b.OutputDirectory)
		fmt.Printf("+ Output minification: %s\n", minificationStatus)
	}

	// --- Decoding ---

	decodeTimeStart := time.Now()
	if !context.NoNonEssentialMessages {
		fmt.Println(terminalStyle("Start decoding...", TerminalTextBold))
	}

	index, err := mkpub.DecodeMainIndex(b.InputDirectory)
	if err != nil {
		return err
	}

	decodeTimeEnd := time.Since(decodeTimeStart)
	if !context.NoNonEssentialMessages {
		fmt.Printf(terminalStyle("Finished decoding!", TerminalTextGreen)+" (%v)\n", decodeTimeEnd)
	}

	// --- Generating ---

	generateTimeStart := time.Now()
	if !context.NoNonEssentialMessages {
		fmt.Println(terminalStyle("Start generating...", TerminalTextBold))
	}

	if err := mkpub.WriteIndexToStaticWebsite(&index, b.OutputDirectory); err != nil {
		return err
	}

	generateTimeEnd := time.Since(generateTimeStart)
	if !context.NoNonEssentialMessages {
		fmt.Printf(terminalStyle("Finished generating!", TerminalTextGreen)+" (%v)\n", generateTimeEnd)
	}

	// --- End build ---
	buildTimeEnd := time.Since(buildTimeStart)
	if !context.NoNonEssentialMessages {
		fmt.Printf(terminalStyle("Finished building!", TerminalTextBold, TerminalTextGreen)+" (%v)\n", buildTimeEnd)
	}

	return nil
}

type VersionCommand struct{}

func (v *VersionCommand) Run() error {
	handleVersionFlag(true)

	return nil
}

type CLI struct {
	Build                  BuildCommand   `cmd:"" help:"Convert source files into a static website and other distributable output formats"`
	Version                VersionCommand `cmd:"" help:"Print application version to the terminal (i.e. standard output) and exit successfully (code 0)"`
	NoNonEssentialMessages bool           `short:"q" help:"Disable printing non-error debug messages (e.g. build progress messages) to the terminal (i.e. standard output)"`
	Plain                  bool           `env:"NO_COLOR" help:"Disable printing colored/styled debug text (i.e. terminal escape codes) to the terminal (i.e. standard output)"`
}

func handleVersionFlag(versionFlag bool) {
	if !versionFlag {
		return
	}

	if len(os.Args) == 0 {
		fmt.Fprintf(os.Stderr, "mkpub error: missing command name")
		os.Exit(1)
	}

	fmt.Printf("%s version %s\n", os.Args[0], Version)
	os.Exit(0)
}

func terminalStyle(s string, codes ...string) string {
	if cli.Plain {
		return s
	}

	return TerminalStyle(s, codes...)
}

func main() {
	context := kong.Parse(&cli)
	if err := context.Run(&Context{
		Plain:                  cli.Plain,
		NoNonEssentialMessages: cli.NoNonEssentialMessages,
	}); err != nil {
		context.FatalIfErrorf(err)
	}
}
