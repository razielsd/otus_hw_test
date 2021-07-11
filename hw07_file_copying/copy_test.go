package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy_EmptySrcPath(t *testing.T) {
	err := Copy("", "path", 100, 0)
	require.Error(t, err)
	require.Equal(t, ErrorEmptySourcePath, err)
}

func TestCopy_EmptyDestPath(t *testing.T) {
	err := Copy("src", "", 100, 0)
	require.Error(t, err)
	require.Equal(t, ErrorEmptyDestinationPath, err)
}

func TestCopy_FileNotFound(t *testing.T) {
	err := Copy("./testdata/file_not_found", "path", 100, 0)
	require.Error(t, err)
	require.Equal(t, ErrorFileNotFound, err)
}

func TestCopy_CopyNonRegular(t *testing.T) {
	err := Copy("./testdata", "/tmp/copyreg.txt", 100, 0)
	require.Error(t, err)
	require.Equal(t, ErrorCopyNonRegular, err)
}

func TestCopy_CopyEmptyFile(t *testing.T) {
	err := Copy("./testdata/emptyfile.txt", "/tmp/copyreg.txt", 100, 0)
	require.Error(t, err)
	require.Equal(t, ErrorEmptySrcFile, err)
}

func TestGetBlockSize_LimitGreaterThanZero(t *testing.T) {
	tests := []struct {
		input    int64
		expected int64
	}{
		{input: 1, expected: 1},
		{input: MinBlockSize, expected: MinBlockSize},
		{input: (BlockCount / 10) * MinBlockSize, expected: MinBlockSize},
		{input: BlockCount * MinBlockSize, expected: MinBlockSize},
		{input: 5 * BlockCount * MinBlockSize, expected: 5 * MinBlockSize},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("size%d", tc.input), func(t *testing.T) {
			size, err := GetBlockSize(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, size)
		})
	}
}

func TestGetBlockSize_LimitLowerOrEqualZero(t *testing.T) {
	_, err := GetBlockSize(0)
	require.Error(t, err)
}
