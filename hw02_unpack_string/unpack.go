package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var b strings.Builder
	var i int
	var prevLetter bool
	var r = []rune(s)

	switch s {
	case "":
		return "", nil
	default:
		for {
			if unicode.IsLetter(r[i]) {
				b.WriteRune(r[i])
				i++
				prevLetter = true
			} else if unicode.IsDigit(r[i]) && i != 0 && prevLetter {
				n, _ := strconv.Atoi(string(r[i]))
				str := b.String()
				last := str[len(str)-1:]
				if n > 0 {
					b.WriteString(strings.Repeat(last, n-1))
				} else if n == 0 {
					s := b.String()
					s = s[:len(s)-1]
					b.Reset()
					b.WriteString(s)
				}
				i++
				prevLetter = false
			} else if r[i] == '\\' {
				i++
				b.WriteRune(r[i])
				i++
				prevLetter = true
			} else {
				return "", ErrInvalidString
			}
			if i == len(s) {
				break
			}
		}
		return b.String(), nil
	}
}
