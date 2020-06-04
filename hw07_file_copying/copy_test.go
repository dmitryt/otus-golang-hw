package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const inTestDir = "testdata"
const outTestDir = ".tmp"

var inFilePath = path.Join(inTestDir, "input.txt")
var outFilePath = path.Join(outTestDir, "out.txt")
var testFilePath = path.Join(outTestDir, "some.txt")

func mkFile(path string, content string) {
	file, err := os.Create(path)
	if err != nil {
		panic("Cannot create test file")
	}
	if _, err := file.WriteString(content); err != nil {
		file.Close()
		panic("Cannot write to test file")
	}
	file.Close()
}

func readFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic("Cannot read test file")
	}
	return string(content)
}

func setup(path string, perm os.FileMode) {
	if err := os.Mkdir(path, perm); err != nil {
		panic("Cannot create test directory")
	}
}

func shutdown(path string) {
	if err := os.RemoveAll(path); err != nil {
		panic("Cannot remove test directory")
	}
}

func TestValidation(t *testing.T) {
	t.Run("should throw validation error, when offset is greater, than file size", func(t *testing.T) {
		setup(outTestDir, os.ModeDir)
		defer shutdown(outTestDir)

		result := Copy(inFilePath, outFilePath, 10000, 0)
		require.Equal(t, WrapError(ErrOffsetExceedsFileSize), result)
	})

	t.Run("should throw validation error, when user tries to copy directory", func(t *testing.T) {
		setup(outTestDir, os.ModeDir)
		defer shutdown(outTestDir)

		result := Copy(inTestDir, outFilePath, 0, 0)
		require.Equal(t, WrapError(ErrSourceIsDirectory), result)
	})

	t.Run("should throw validation error, when user tries to copy directory", func(t *testing.T) {
		setup(outTestDir, os.ModeDir)
		defer shutdown(outTestDir)

		result := Copy("/dev/random", outFilePath, 0, 0)
		require.Equal(t, WrapError(ErrUnsupportedFile), result)
	})

	t.Run("should throw validation error, when user tries to copy directory", func(t *testing.T) {
		setup(outTestDir, os.ModeDir)
		defer shutdown(outTestDir)

		result := Copy("/dev/random", outFilePath, 0, 0)
		require.Equal(t, WrapError(ErrUnsupportedFile), result)
	})

	t.Run("should throw validation error, when user tries to write to directory", func(t *testing.T) {
		setup(outTestDir, os.ModeDir)
		defer shutdown(outTestDir)

		var pathError *os.PathError
		result := Copy(inFilePath, outTestDir, 0, 0)
		require.True(t, errors.As(result, &pathError))
	})

	t.Run("should throw validation error, when user tries to write to readonly directory", func(t *testing.T) {
		readonlyDir := ".readonlydir"
		setup(readonlyDir, 0444)
		defer shutdown(readonlyDir)

		var pathError *os.PathError
		result := Copy(inFilePath, path.Join(readonlyDir, "some.txt"), 0, 0)
		require.True(t, errors.As(result, &pathError))
	})
}

func TestCopy(t *testing.T) {
	t.Run("should copy all file", func(t *testing.T) {
		fileContent := "some content"
		setup(outTestDir, 0755)
		defer shutdown(outTestDir)

		mkFile(testFilePath, fileContent)

		require.Nil(t, Copy(testFilePath, outFilePath, 0, 0))
		require.Equal(t, fileContent, readFile(outFilePath))
	})

	t.Run("should copy file with offset", func(t *testing.T) {
		fileContent := "some content"
		setup(outTestDir, 0755)
		defer shutdown(outTestDir)
		mkFile(testFilePath, fileContent)

		require.Nil(t, Copy(testFilePath, outFilePath, 5, 0))
		require.Equal(t, "content", readFile(outFilePath))
	})

	t.Run("should copy file with limit", func(t *testing.T) {
		fileContent := "some content"
		setup(outTestDir, 0755)
		defer shutdown(outTestDir)
		mkFile(testFilePath, fileContent)

		require.Nil(t, Copy(testFilePath, outFilePath, 0, 5))
		require.Equal(t, "some ", readFile(outFilePath))
	})

	t.Run("should copy file with offset/limit", func(t *testing.T) {
		fileContent := "some content"
		setup(outTestDir, 0755)
		defer shutdown(outTestDir)
		mkFile(testFilePath, fileContent)

		require.Nil(t, Copy(testFilePath, outFilePath, 2, 5))
		require.Equal(t, "me co", readFile(outFilePath))
	})
}
