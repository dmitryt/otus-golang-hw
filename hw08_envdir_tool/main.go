package main

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrReadDir     = errors.New("cannot read directory")
	ErrExecCommand = errors.New("the command cannot be executed")
)

const defaultErrCode = 111

func wrapError(err error) error {
	return fmt.Errorf("envdir: %w", err)
}

func main() {
	dirPath, args := os.Args[1], os.Args[2:]
	env, err := ReadDir(dirPath)
	if err != nil {
		fmt.Println(wrapError(ErrReadDir))
		os.Exit(defaultErrCode)
	}
	os.Exit(RunCmd(args, env))
}
