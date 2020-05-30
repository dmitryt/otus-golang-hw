package main

import (
	"strings"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cmdExec := exec.Command(cmd[0], cmd[1:]...)
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
	err := cmdExec.Start()
	if err != nil {
		fmt.Println(wrapError(err))
		return 111
	}
	cmdExec.Wait()
	return cmdExec.ProcessState.ExitCode()
}
