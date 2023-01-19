package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

//easyjson:json
type User struct {
	ID       int    `json:"-"`
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	var d strings.Builder
	d.WriteRune('.')
	d.WriteString(domain)

	fileScanner := bufio.NewScanner(r)
	var user User
	result := make(DomainStat)
	dm := d.String()
	for fileScanner.Scan() {
		if err := easyjson.Unmarshal(fileScanner.Bytes(), &user); err != nil {
			return nil, err
		}
		if strings.Contains(user.Email, dm) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
