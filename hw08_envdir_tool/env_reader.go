package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	if strings.Contains(dir, " ") {
		return nil, errors.New("folder exists whitespace")
	}
	e := make(Environment)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, fi := range files {
		f, err := os.Open(path.Join(dir, fi.Name()))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)
		scanner.Scan()
		str := scanner.Text()

		s := strings.ReplaceAll(str, "\u0000", "\n")
		s = trimEndValidate(s)

		if str == "" {
			e[fi.Name()] = EnvValue{Value: "", NeedRemove: true}
		} else {
			e[fi.Name()] = EnvValue{Value: s, NeedRemove: false}
		}
		f.Close()
	}
	return e, nil
}

// Search last validate symbol.
func trimEndValidate(s string) string {
	r := []rune(s)
	for i := len(r) - 1; i > 0; i-- {
		if r[i] != ' ' && r[i] != '\t' {
			return string(r[:i+1])
		}
	}
	return s
}
