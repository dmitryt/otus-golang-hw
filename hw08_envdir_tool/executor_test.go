package main

import (
	"os"
	"testing"

	"github.com/kami-zh/go-capturer"
	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("should run command and return exit codes correctly", func(t *testing.T) {
		env := Environment{}
		require.Equal(t, 0, RunCmd([]string{"ls", "testdata"}, env))
		require.Equal(t, 1, RunCmd([]string{"cat", "doesnotexist.txt"}, env))
		require.Equal(t, 111, RunCmd([]string{"non-existing-command"}, env))
	})

	t.Run("should capture stdout/stderr correctly", func(t *testing.T) {
		env := Environment{}
		outData := capturer.CaptureStdout(func() {
			RunCmd([]string{"ls", "."}, env)
		})
		require.Contains(t, outData, "main.go")
		errData := capturer.CaptureStderr(func() {
			RunCmd([]string{"cat", "doesnotexist.txt"}, env)
		})
		require.Contains(t, errData, "No such file or directory")
	})

	t.Run("should apply provided env vars", func(t *testing.T) {
		env := Environment{"HELLO": "WORLD"}
		os.Setenv("PREDEFINED", "SOME")
		output := capturer.CaptureStdout(func() {
			RunCmd([]string{"env"}, env)
		})
		require.Contains(t, output, "PREDEFINED=SOME")
		require.Contains(t, output, "HELLO=WORLD")
	})

	t.Run("should update existing env vars", func(t *testing.T) {
		env := Environment{"HELLO": "WORLD"}
		os.Setenv("HELLO", "VALUE")
		output := capturer.CaptureStdout(func() {
			RunCmd([]string{"env"}, env)
		})
		require.Contains(t, output, "HELLO=WORLD")
	})

	t.Run("should unset existing env vars correctly", func(t *testing.T) {
		env := Environment{"HELLO": ""}
		os.Setenv("HELLO", "VALUE")
		output := capturer.CaptureStdout(func() {
			RunCmd([]string{"env"}, env)
		})
		require.NotContains(t, output, "HELLO")
	})

	t.Run("should unset env var and skip setting value, if it contains '='", func(t *testing.T) {
		env := Environment{"HELLO": "SOME=VALUE"}
		os.Setenv("HELLO", "VALUE")
		output := capturer.CaptureStdout(func() {
			RunCmd([]string{"env"}, env)
		})
		require.NotContains(t, output, "HELLO")
	})
}
