package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"

	"github.com/cheggaaa/pb/v3"
)

const (
	BlockCount   = 100
	MinBlockSize = 10
)

var (
	ErrorEmptySourcePath      = errors.New("empty source path")
	ErrorEmptyDestinationPath = errors.New("empty destination path")
	ErrorFileNotFound         = errors.New("file not found")
	ErrorUnableGetStat        = errors.New("unable get file stat")
	ErrorCopyNonRegular       = errors.New("unable to copy non-regular file")
	ErrOffsetExceedsFileSize  = errors.New("offset greater than limit")
	ErrorCreateTmpFile        = errors.New("unable to create tmp file")
	ErrorCopyFile             = errors.New("error due reading file")
	ErrorEmptySrcFile         = errors.New("empty source file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if len(fromPath) == 0 {
		return ErrorEmptySourcePath
	}

	if len(toPath) == 0 {
		return ErrorEmptyDestinationPath
	}

	srcFile, err := os.Open(fromPath)
	if os.IsNotExist(err) {
		return ErrorFileNotFound
	} else if err != nil {
		return err
	}

	defer srcFile.Close()

	si, err := srcFile.Stat()
	if err != nil {
		return ErrorUnableGetStat
	}

	if !si.Mode().IsRegular() {
		return ErrorCopyNonRegular
	}

	if si.Size() == 0 {
		return ErrorEmptySrcFile
	}

	if si.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 {
		limit = si.Size()
	}

	if offset+limit > si.Size() {
		limit = si.Size() - offset
	}

	destFile, err := ioutil.TempFile("/tmp", "cp")
	if err != nil {
		return ErrorCreateTmpFile
	}

	blockSize, err := GetBlockSize(limit)
	if err != nil {
		return err
	}
	var cpsize int64 = 0
	bar := pb.Full.Start64(limit)
	for cpsize < limit {
		if cpsize+blockSize > limit {
			blockSize = limit - cpsize
		}
		buff := make([]byte, blockSize)
		_, err = srcFile.ReadAt(buff, offset)
		if err != nil {
			return ErrorCopyFile
		}
		barReader := bar.NewProxyReader(bytes.NewReader(buff))
		written, err := io.CopyN(destFile, barReader, blockSize)
		if err != nil {
			return err
		}
		cpsize += written
		offset += blockSize
	}
	bar.Finish()

	if err = destFile.Close(); err != nil {
		return err
	}
	if err = os.Rename(destFile.Name(), toPath); err != nil {
		return err
	}

	return nil
}

func GetBlockSize(limit int64) (int64, error) {
	if limit < 1 {
		return 0, fmt.Errorf("logic error, read limit(%d) is lower than 1", limit)
	}
	size := math.Ceil(float64(limit) / float64(BlockCount))
	si := int64(size)
	if si < MinBlockSize {
		si = MinBlockSize
	}
	if si > limit {
		si = limit
	}
	return si, nil
}
