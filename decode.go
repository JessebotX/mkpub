package mkpub

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

const (
	IndexConfigName = "mkpub.yml"
	BookConfigName  = "mkpub-book.yml"
)

type RenderIndex struct {
	Index

	InputPath string
}

func (i *RenderIndex) InitDefaults(inputPath string) error {
	i.InputPath = inputPath

	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return err
	}

	i.Title = filepath.Base(absInputPath)
	i.LayoutsDirectory = filepath.Join(inputPath, "layout")

	return nil
}

func DecodeIndex(inputPath string) (RenderIndex, error) {
	var index RenderIndex
	if err := index.InitDefaults(inputPath); err != nil {
		return index, fmt.Errorf("index: failed to initialize Index: %w", err)
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

	return index, nil
}

func mapToStruct(m map[string]any, s any) error {
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, s); err != nil {
		return err
	}

	return nil
}
