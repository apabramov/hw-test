package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	c := exec.Command(cmd[0], cmd[1:]...)

	//	fmt.Println(c.String(), env)
	for i, v := range env {
		os.Unsetenv(i)
		fmt.Printf("%v=%v", i, v.Value)
		if !v.NeedRemove {
			c.Env = append(os.Environ(), fmt.Sprintf("%s=%s", i, v.Value))
			//	fmt.Printf("%s=%s", i, v.Value)
		}
	}

	//stdin, err := c.StdinPipe()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer stdin.Close()
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Start(); err != nil {
		fmt.Println(err)
	}

	// io.WriteString(stdin, "4\n")
	if err := c.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return exiterr.ExitCode()
		} else {
			fmt.Println(err)
		}
	}
	return
}
