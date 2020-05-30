package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cmdExec := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	envList := []string{}
	for key, value := range env {
		if _, ok := os.LookupEnv(key); ok {
			os.Unsetenv(key)
		}
		if strings.Contains(value, "=") {
			continue
		}
		if value != "" {
			envList = append(envList, fmt.Sprintf("%s=%s", key, value))
		}
	}
	cmdExec.Env = append(os.Environ(), envList...)
	cmdExec.Stdin = os.Stdin
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	err := cmdExec.Run()
	var eerr *exec.Error
	if errors.As(err, &eerr) {
		fmt.Println(wrapError(ErrExecCommand))
		return defaultErrCode
	}
	return cmdExec.ProcessState.ExitCode()
}
