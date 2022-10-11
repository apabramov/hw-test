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
			case unicode.IsLetter(r[i]) || r[i] == ' ':
				b.WriteRune(r[i])
				i++
				prevLetter = true
			case unicode.IsDigit(r[i]) && i != 0 && prevLetter:
				n, _ := strconv.Atoi(string(r[i]))
				s := []rune(b.String())
				switch {
				case n > 0:
					b.WriteString(strings.Repeat(string(s[len(s)-1:]), n-1))
				case n == 0:
					s = s[:len(s)-1]
					b.Reset()
					b.WriteString(string(s))
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
			if i == len([]rune(s)) {
				break
			}
		}
		return b.String(), nil
	}
}
