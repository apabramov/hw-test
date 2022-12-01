package main

import (
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	var l int64
	inFile, err := os.Open(fromPath)
	if err != nil {
		fmt.Println(err.Error())
	}

	outFile, err := os.Create(toPath)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		inFile.Close()
		outFile.Close()
	}()

	_, err = inFile.Seek(offset, io.SeekStart)
	if err != nil {
		fmt.Println(err.Error())
	}

	f, err := inFile.Stat()
	if err != nil {
		fmt.Println(err.Error())
	}

	if limit == 0 || limit > f.Size() {
		l = f.Size()
	} else {
		l = limit
	}

	bar := pb.Full.Start64(l)
	b := bar.NewProxyWriter(outFile)

	_, er := io.CopyN(b, inFile, l)
	if er != nil {
		fmt.Println(er.Error())
	}

	bar.Finish()

	outFile.Close()
	inFile.Close()

	return nil
}
