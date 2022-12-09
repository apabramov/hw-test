package main

import (
	"bufio"
	"bytes"
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
	e := make(Environment)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, fi := range files {
		err := func() error {
			f, err := os.Open(path.Join(dir, fi.Name()))
			if err != nil {
				return err
			}
			defer f.Close()
			//r := bufio.NewReader(f)
			//str, er := r.ReadString('\r')
			scanner := bufio.NewScanner(f)
			scanner.Split(bufio.ScanLines)
			scanner.Scan()
			str := scanner.Text()
			//fmt.Println("-", str, "-")
			s := bytes.Replace([]byte(str), []byte{0}, []byte{10}, -1)
			//_ = s
			//fmt.Println("env", fi.Name(), str, strings.TrimRight(string(s), "\n"))

			ss := strings.TrimRight(string(s), "\r")
			ss = strings.TrimRight(ss, " ")
			ss = strings.TrimRight(ss, "\t")
			if str == "" {
				e[fi.Name()] = EnvValue{Value: "", NeedRemove: true}
			} else {
				e[fi.Name()] = EnvValue{Value: ss, NeedRemove: false}
			}
			return nil
		}()
		if err != nil {
			return nil, err
		}
	}
	return e, nil
}
