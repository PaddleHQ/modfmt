package modfmt_test

import (
	"os"
	"testing"

	"github.com/PaddleHQ/modfmt/pkg/modfmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/mod/modfile"
)

const (
	testFileName        = "testdata/go.mod"
	updatedTestFileName = "testdata/updated_go.mod"
)

func TestMergeRequires(t *testing.T) {
	// fmt the go.mod file
	updatedContents, err := modfmt.MergeRequires(testFileName)
	require.NoError(t, err)

	// create a new file with the updated contents
	newFile, err := os.Create(updatedTestFileName)
	require.NoError(t, err)

	defer func() {
		err := newFile.Close()
		assert.NoError(t, err)
	}()

	_, err = newFile.Write(updatedContents)
	require.NoError(t, err)

	// parse the old go.mod file
	oldContents, err := os.ReadFile(testFileName)
	require.NoError(t, err)

	oldmod, err := modfile.ParseLax(testFileName, oldContents, nil)
	require.NoError(t, err)

	// parse the updated go.mod file
	newContents, err := os.ReadFile(updatedTestFileName)
	require.NoError(t, err)

	newmod, err := modfile.ParseLax(updatedTestFileName, newContents, nil)
	require.NoError(t, err)

	// Check if both mod reqs have the same length
	assert.Equal(t, len(oldmod.Require), len(newmod.Require))

	// Check if both mod reqs have the same content
	for _, req := range oldmod.Require {
		found := false
		for _, newReq := range newmod.Require {
			if req.Mod.Path == newReq.Mod.Path && req.Mod.Version == newReq.Mod.Version {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Require mismatch: %s@%s not found in updated go.mod", req.Mod.Path, req.Mod.Version)
		}
	}
}
