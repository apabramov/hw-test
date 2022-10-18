package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var re = regexp.MustCompile(`[\p{L}\d_[^-]+`)

func Top10(s string) []string {
	str := make([]string, 0)
	m := make(map[string]int)

	if s == "" {
		return nil
	}

	for _, v := range re.FindAllString(s, -1) {
		l := strings.ToLower(v)
		if _, found := m[l]; found && v != "-" {
			m[l]++
		} else {
			m[l] = 1
		}
	}
	for i := range m {
		str = append(str, i)
	}
	sort.Slice(str, func(i, j int) bool {
		return (str[i] < str[j] && m[str[i]] == m[str[j]]) || (m[str[i]] > m[str[j]])
	})

	return str[:10]
}
