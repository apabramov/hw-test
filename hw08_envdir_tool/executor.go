package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	com := cmd[0]
	c := exec.Command(com, cmd[1:]...)

	for i, v := range env {
		os.Unsetenv(i)
		if !v.NeedRemove {
			err := os.Setenv(i, v.Value)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	c.Env = os.Environ()

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Start(); err != nil {
		fmt.Println(err)
	}

	if err := c.Wait(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
	}
	return 0
}
