package main

import (
	"bufio"
	"bytes"
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
		func() {
			f, err := os.Open(path.Join(dir, fi.Name()))
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			scanner.Split(bufio.ScanLines)
			scanner.Scan()
			str := scanner.Text()

			s := bytes.ReplaceAll([]byte(str), []byte{0}, []byte{10})

			val := strings.TrimRight(string(s), "\r")
			val = strings.TrimRight(val, "\t")
			val = strings.TrimRight(val, " ")
			val = strings.TrimRight(val, "\t")
			if str == "" {
				e[fi.Name()] = EnvValue{Value: "", NeedRemove: true}
			} else {
				e[fi.Name()] = EnvValue{Value: val, NeedRemove: false}
			}
		}()
	}
	return e, nil
}
