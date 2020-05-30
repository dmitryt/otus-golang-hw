package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const testDir = ".tmp"

var outFilePath = path.Join(testDir, "out.txt")

func setup() {
	if err := os.Mkdir(testDir, 0755); err != nil {
		panic("Cannot create test directory")
	}
}

func shutdown() {
	if err := os.RemoveAll(testDir); err != nil {
		panic("Cannot remove test directory")
	}
}

func prepareFile(fileName string, content string) {
	file, err := os.Create(path.Join(testDir, fileName))
	if err != nil {
		panic("Cannot create test file")
	}
	defer file.Close()
	if _, err := file.WriteString(content); err != nil {
		panic("Cannot write to test file")
	}
}

func prepareFileWithPerm(fileName string, content string, perm os.FileMode) {
	prepareFile(fileName, content)
	if err := os.Chmod(path.Join(testDir, fileName), perm); err != nil {
		panic("Cannot change file permission")
	}
}

func TestReadDir(t *testing.T) {
	t.Run("should read data from files correctly", func(t *testing.T) {
		setup()
		defer shutdown()
		prepareFile("HELLO", "WORLD")
		prepareFile("VAR_WITH_SPACES", "  VALUE 	 ")
		prepareFile("VAR_EMPTY", "")
		prepareFile("VAR_WITH_TERMINAL_NULLS", "V\x00A\000L\u0000UE")
		prepareFile("VAR_WITH_MULTILINE", "VALUE\nSECOND LINE\nTHIRD LINE")

		result, err := ReadDir(testDir)

		require.Nil(t, err)
		require.Equal(t, Environment{
			"HELLO":                   "WORLD",
			"VAR_WITH_SPACES":         "  VALUE",
			"VAR_EMPTY":               "",
			"VAR_WITH_TERMINAL_NULLS": "V\nA\nL\nUE",
			"VAR_WITH_MULTILINE":      "VALUE",
		}, result)
	})

	t.Run("should read data from files without read permissions", func(t *testing.T) {
		setup()
		defer shutdown()
		prepareFileWithPerm("HELLO", "WORLD", 0644)
		prepareFileWithPerm("SOME_VAR", "VALUE", 0111)
		prepareFileWithPerm("ANOTHER_VAR", "VALUE", 0222)

		result, err := ReadDir(testDir)

		require.Nil(t, err)
		require.Equal(t, Environment{
			"HELLO": "WORLD",
		}, result)
	})
}
