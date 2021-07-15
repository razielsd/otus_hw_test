package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy_EmptySrcPath(t *testing.T) {
	err := Copy("", "path", 100, 0)
	require.ErrorIs(t, ErrEmptySourcePath, err)
}

func TestCopy_EmptyDestPath(t *testing.T) {
	err := Copy("src", "", 100, 0)
	require.ErrorIs(t, ErrEmptyDestinationPath, err)
}

func TestCopy_FileNotFound(t *testing.T) {
	err := Copy("./testdata/file_not_found", "path", 100, 0)
	require.ErrorIs(t, ErrFileNotFound, err)
}

func TestCopy_CopyNonRegular(t *testing.T) {
	err := Copy("./testdata", "/tmp/copyreg.txt", 100, 0)
	require.ErrorIs(t, ErrCopyNonRegular, err)
}

func TestCopy_CopyEmptyFile(t *testing.T) {
	err := Copy("./testdata/emptyfile.txt", "/tmp/copyreg.txt", 100, 0)
	require.ErrorIs(t, ErrEmptySrcFile, err)
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
			expData, err := ioutil.ReadFile(tc.expected)
			require.NoError(t, err)
			destData, err := ioutil.ReadFile(dest)
			require.NoError(t, err)
			require.Equal(t, expData, destData, "Files not equal")
		})
	}
}
