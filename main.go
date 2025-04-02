// Package main is responsible for handling CLI entrypoint of modfmt
package main

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/PaddleHQ/modfmt/pkg/modfmt"
)

const gomodName = "go.mod"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	updatedContents, err := modfmt.MergeRequires(gomodName)
	if err != nil {
		return fmt.Errorf("failed to merge requires: %w", err)
	}

	// get arguments
	if !slices.Contains(os.Args, "--replace") {
		// print updated contents to stdout
		//nolint:forbidigo // This is a CLI tool
		fmt.Println(string(updatedContents))
		return nil
	}

	// write updated contents to go.mod
	info, err := os.Stat(gomodName)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	if err = os.WriteFile(gomodName, updatedContents, info.Mode()); err != nil {
		return fmt.Errorf("failed to write updated go.mod: %w", err)
	}

	return nil
}
