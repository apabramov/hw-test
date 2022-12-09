package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	c := exec.Command(cmd[0], cmd[1:]...)

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
		if exiterr, ok := err.(*exec.ExitError); ok {
			return exiterr.ExitCode()
		} else {
			fmt.Println(err)
		}
	}
	return 0
}
