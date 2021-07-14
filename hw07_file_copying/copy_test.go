package main

import (
	"fmt"
	"testing"

	"github.com/hlubek/readercomp"
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

/*
./go-cp -from testdata/input.txt -to out.txt
./go-cp -from testdata/input.txt -to out.txt -limit 10
./go-cp -from testdata/input.txt -to out.txt -limit 1000
./go-cp -from testdata/input.txt -to out.txt -limit 10000
./go-cp -from testdata/input.txt -to out.txt -offset 100 -limit 1000
./go-cp -from testdata/input.txt -to out.txt -offset 6000 -limit 1000
*/

func TestCopy_ValidParams_SuccesCopy(t *testing.T) {
	tests := []struct {
		src      string
		dest     string
		limit    int64
		offset   int64
		expected string
	}{
		{
			src:      "./testdata/input.txt",
			dest:     "out.txt",
			limit:    0,
			offset:   0,
			expected: "testdata/out_offset0_limit0.txt",
		},
		{
			src:      "./testdata/input.txt",
			dest:     "out.txt",
			limit:    10,
			offset:   0,
			expected: "testdata/out_offset0_limit10.txt",
		},
		{
			src:      "./testdata/input.txt",
			dest:     "out.txt",
			limit:    1000,
			offset:   0,
			expected: "testdata/out_offset0_limit1000.txt",
		},
		{
			src:      "./testdata/input.txt",
			dest:     "out.txt",
			limit:    10000,
			offset:   0,
			expected: "testdata/out_offset0_limit10000.txt",
		},
		{
			src:      "./testdata/input.txt",
			dest:     "out.txt",
			limit:    1000,
			offset:   100,
			expected: "testdata/out_offset100_limit1000.txt",
		},
		{
			src:      "./testdata/input.txt",
			dest:     "out.txt",
			limit:    1000,
			offset:   6000,
			expected: "testdata/out_offset6000_limit1000.txt",
		},
	}
	for _, tc := range tests {
		tc := tc
		name := fmt.Sprintf("cp_limit_%d_offset_%d", tc.limit, tc.offset)
		t.Run(name, func(t *testing.T) {
			dest := t.TempDir() + "/" + tc.dest + ".txt"
			err := Copy(tc.src, dest, tc.offset, tc.limit)
			require.NoError(t, err)
			require.FileExists(t, dest)
			ok, err := readercomp.FilesEqual(tc.expected, dest)
			require.NoError(t, err, "error compare files")
			require.True(t, ok, "copied file is not expected")
		})
	}
}
