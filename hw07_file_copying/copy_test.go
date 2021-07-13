package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy_EmptySrcPath(t *testing.T) {
	err := Copy("", "path", 100, 0)
	require.ErrorIs(t, ErrorEmptySourcePath, err)
}

func TestCopy_EmptyDestPath(t *testing.T) {
	err := Copy("src", "", 100, 0)
	require.ErrorIs(t, ErrorEmptyDestinationPath, err)
}

func TestCopy_FileNotFound(t *testing.T) {
	err := Copy("./testdata/file_not_found", "path", 100, 0)
	require.ErrorIs(t, ErrorFileNotFound, err)
}

func TestCopy_CopyNonRegular(t *testing.T) {
	err := Copy("./testdata", "/tmp/copyreg.txt", 100, 0)
	require.ErrorIs(t, ErrorCopyNonRegular, err)
}

func TestCopy_CopyEmptyFile(t *testing.T) {
	err := Copy("./testdata/emptyfile.txt", "/tmp/copyreg.txt", 100, 0)
	require.ErrorIs(t, ErrorEmptySrcFile, err)
}
