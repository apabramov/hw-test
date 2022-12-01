package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	var l int64

	inFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	inF, err := inFile.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}

	if offset > inF.Size() {
		return ErrOffsetExceedsFileSize
	}

	outFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = inFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	if limit == 0 || limit > inF.Size() {
		l = inF.Size()
	} else {
		l = limit
	}

	bar := pb.Full.Start64(l)
	b := bar.NewProxyWriter(outFile)

	if _, err := io.CopyN(b, inFile, l); err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	bar.Finish()
	return nil
}
