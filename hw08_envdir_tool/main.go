package main

import (
	"fmt"
	"os"
	"errors"
)

var (
	ErrReadDir = errors.New("Cannot read directory")
)

func wrapError(err error) error {
	return fmt.Errorf("envdir: %w", err)
}

func main() {
	dirPath, args := os.Args[1], os.Args[2:]
	env, err := ReadDir(dirPath)
	if err != nil {
		fmt.Println(wrapError(ErrReadDir))
		os.Exit(111)
	}
	os.Exit(RunCmd(args, env))
}
