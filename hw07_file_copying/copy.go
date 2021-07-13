package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/cheggaaa/pb/v3"
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
		return fmt.Errorf("cannot open file: %w", err)
	}

	defer srcFile.Close()

	si, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorUnableGetStat, err)
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

	if err != nil {
		return err
	}
	bar := pb.Full.Start64(limit)
	defer bar.Finish()
	buff := make([]byte, limit)
	_, err = srcFile.ReadAt(buff, offset)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorCopyFile, err)
	}
	barReader := bar.NewProxyReader(bytes.NewReader(buff))
	_, err = io.CopyN(destFile, barReader, limit)

	if err != nil {
		return err
	}

	if err = destFile.Close(); err != nil {
		return err
	}
	if err = os.Rename(destFile.Name(), toPath); err != nil {
		return err
	}

	return nil
}
