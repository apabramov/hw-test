package main

import (
	"fmt"
	"os"
)

func main() {
	fl := os.Args
	env, err := ReadDir(fl[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Exit(RunCmd(fl[2:], env))
}
