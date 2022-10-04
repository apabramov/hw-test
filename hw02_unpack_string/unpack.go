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
	r := []rune(s)

	switch s {
	case "":
		return "", nil
	default:
		for {
			switch {
			case unicode.IsLetter(r[i]):
				b.WriteRune(r[i])
				i++
				prevLetter = true
			case unicode.IsDigit(r[i]) && i != 0 && prevLetter:
				n, _ := strconv.Atoi(string(r[i]))
				str := b.String()
				last := str[len(str)-1:]
				switch {
				case n > 0:
					b.WriteString(strings.Repeat(last, n-1))
				case n == 0:
					s := b.String()
					s = s[:len(s)-1]
					b.Reset()
					b.WriteString(s)
				}
				i++
				prevLetter = false
			case r[i] == '\\':
				i++
				b.WriteRune(r[i])
				i++
				prevLetter = true
			default:
				return "", ErrInvalidString
			}
			if i == len(s) {
				break
			}
		}
		return b.String(), nil
	}
}
