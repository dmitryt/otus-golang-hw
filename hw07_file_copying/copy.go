package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSourceIsDirectory     = errors.New("cannot copy a directory")
)

type File struct {
	ref  *os.File
	stat os.FileInfo
}

const KB = 1 << 10

func WrapError(err error) error {
	return fmt.Errorf("go-cp: %w", err)
}

func getSource(path string) (*File, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if stat.Mode().IsDir() {
		return nil, ErrSourceIsDirectory
	}
	if !stat.Mode().IsRegular() {
		return nil, ErrUnsupportedFile
	}
	return &File{ref: file, stat: stat}, nil
}

func writeContent(src io.Reader, toPath string, limit int64, onCopy func(int64)) error {
	// Init Tmp File
	tmpFile, err := ioutil.TempFile("", "tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	bufSize := int64(1 * KB)
	for copied := int64(0); copied < limit; {
		if copied+bufSize > limit {
			bufSize = limit - copied
		}
		written, err := io.CopyN(tmpFile, src, bufSize)
		copied += written
		onCopy(copied)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	if err := os.Rename(tmpFile.Name(), toPath); err != nil {
		return err
	}
	return nil
}

func Copy(fromPath string, toPath string, offset, limit int64) error {
	src, err := getSource(fromPath)
	if err != nil {
		return WrapError(err)
	}
	defer src.ref.Close()
	if offset > src.stat.Size() {
		return WrapError(ErrOffsetExceedsFileSize)
	}
	if limit == 0 || limit > src.stat.Size() {
		limit = src.stat.Size()
	} else if limit+offset > src.stat.Size() {
		limit = src.stat.Size() - offset
	}
	if _, err := src.ref.Seek(offset, 0); err != nil {
		return WrapError(err)
	}
	// Check, whether it's possible to create file in provided location
	if _, err := os.Create(toPath); err != nil {
		return WrapError(err)
	}
	// Init Progress Bar
	bar := pb.Full.Start64(limit)
	defer bar.Finish()
	onCopy := func(written int64) {
		bar.SetCurrent(written)
		time.Sleep(time.Millisecond)
	}
	if err = writeContent(src.ref, toPath, limit, onCopy); err != nil {
		return WrapError(err)
	}
	return nil
}
