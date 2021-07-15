package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrEmptySourcePath       = errors.New("empty source path")
	ErrEmptyDestinationPath  = errors.New("empty destination path")
	ErrFileNotFound          = errors.New("file not found")
	ErrUnableGetStat         = errors.New("unable get file stat")
	ErrCopyNonRegular        = errors.New("unable to copy non-regular file")
	ErrOffsetExceedsFileSize = errors.New("offset greater than limit")
	ErrCreateTmpFile         = errors.New("unable to create tmp file")
	ErrCopyFile              = errors.New("error due reading file")
	ErrEmptySrcFile          = errors.New("empty source file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if len(fromPath) == 0 {
		return ErrEmptySourcePath
	}

	if len(toPath) == 0 {
		return ErrEmptyDestinationPath
	}

	srcFile, err := os.Open(fromPath)
	if os.IsNotExist(err) {
		return ErrFileNotFound
	} else if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}

	defer srcFile.Close()

	si, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnableGetStat, err)
	}

	if !si.Mode().IsRegular() {
		return ErrCopyNonRegular
	}

	if si.Size() == 0 {
		return ErrEmptySrcFile
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
		return ErrCreateTmpFile
	}

	bar := pb.Full.Start64(limit)
	defer bar.Finish()
	_, err = srcFile.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCopyFile, err)
	}

	barReader := bar.NewProxyReader(io.LimitReader(srcFile, limit))
	_, err = io.CopyN(destFile, barReader, limit)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrCopyFile, err)
	}

	if err = destFile.Close(); err != nil {
		return fmt.Errorf("%w: %v", ErrCopyFile, err)
	}
	if err = os.Rename(destFile.Name(), toPath); err != nil {
		return fmt.Errorf("%w: %v", ErrCopyFile, err)
	}

	return nil
}
