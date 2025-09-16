package main

import (
	"fmt"
)

type VersionCommand struct{}

func (VersionCommand) Run() error {
	fmt.Printf("pub version %s\n", Version)
	return nil
}
