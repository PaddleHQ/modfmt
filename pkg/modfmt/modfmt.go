// Package modfmt provides methods for formatting go.mod files
package modfmt

import (
	"fmt"
	"os"

	"golang.org/x/mod/modfile"
)

// MergeRequires takes a go.mod file and merges all the require blocks into one.
func MergeRequires(goModFilename string) ([]byte, error) {
	//nolint:gosec // this is just used to read the go.mod file
	contents, err := os.ReadFile(goModFilename)
	if err != nil {
		return nil, fmt.Errorf("[modfmt] failed to read go.mod file: %w", err)
	}

	mod, err := modfile.ParseLax(goModFilename, contents, nil)
	if err != nil {
		return nil, fmt.Errorf("[modfmt] failed to parse go.mod file: %w", err)
	}

	if err := mergeRequires(mod); err != nil {
		return nil, fmt.Errorf("[modfmt] failed to merge requires: %w", err)
	}

	updatedContents, err := mod.Format()
	if err != nil {
		return nil, fmt.Errorf("[modfmt] failed to format go.mod file: %w", err)
	}

	return updatedContents, nil
}

func mergeRequires(mod *modfile.File) (err error) {
	defer func() {
		if r := recover(); r != nil {
			possibleErr, ok := r.(error)
			if ok {
				err = possibleErr
			}
		}
	}()

	allRequires := make([]modfile.Require, len(mod.Require))
	for i, reqs := range mod.Require {
		// Save all the requires to a new slice
		allRequires[i] = *reqs

		// while removing them from the original slice
		if err := mod.DropRequire(reqs.Mod.Path); err != nil {
			return fmt.Errorf("failed to drop require %s: %w", reqs.Mod.Path, err)
		}
	}

	// Cleanup the modfile
	// This removes the empty require blocks
	mod.Cleanup()

	// Add the requires back to the modfile
	for _, reqs := range allRequires {
		mod.AddNewRequire(reqs.Mod.Path, reqs.Mod.Version, reqs.Indirect)
	}
	mod.Cleanup()

	// Sort the require blocks
	mod.SortBlocks()
	mod.Cleanup()

	// Set the require blocks to separate indirects
	mod.SetRequireSeparateIndirect(mod.Require)
	mod.Cleanup()

	return
}
