package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

type lightUser struct {
	Email string
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	cnt := make(DomainStat)
	var user *lightUser
	reader := bufio.NewReader(r)
	for {
		line, _, _ := reader.ReadLine()
		if line == nil {
			break
		}
		if err := json.Unmarshal(line, &user); err != nil {
			fmt.Println("ERR", err)
			continue
		}
		matched := strings.Contains(user.Email, "."+domain)
		if matched {
			domain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			cnt[domain]++
			user.Email = ""
		}
	}
	return cnt, nil
}
