package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "source file to read from")
	flag.StringVar(&to, "to", "", "destination file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	err := Copy(from, to, offset, limit)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmptySourcePath):
			flag.Usage()
		case errors.Is(err, ErrEmptyDestinationPath):
			flag.Usage()
		case errors.Is(err, ErrFileNotFound):
			fmt.Printf("Error: File not exists: %s\n", from)
		case errors.Is(err, ErrUnableGetStat):
			fmt.Printf("Error: Unable get file stat: %s\n", from)
		case errors.Is(err, ErrCopyNonRegular):
			fmt.Printf("Error: Unable to copy non-regular file: %s\n", from)
		case errors.Is(err, ErrOffsetExceedsFileSize):
			fmt.Printf("Error: Offset if greater than file size\n")
		case errors.Is(err, ErrCreateTmpFile):
			fmt.Printf("Error: Unable to create tmp file\n")
		default:
			fmt.Printf("Error: %s\n", err)
		}
		os.Exit(1)
	}
}
