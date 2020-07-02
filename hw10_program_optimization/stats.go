package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"io"
	"strings"
	"sync"

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

type safeCounter struct {
	v   DomainStat
	mux sync.Mutex
}

type lightUser struct {
	Email string
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	var pool = sync.Pool{
		New: func() interface{} { return new(lightUser) },
	}
	cnt := safeCounter{v: make(DomainStat)}

	scanner := bufio.NewScanner(r)
	var wg sync.WaitGroup
	for scanner.Scan() {
		wg.Add(1)
		go func(content []byte) {
			defer wg.Done()
			user := pool.Get().(*lightUser)
			defer pool.Put(user)
			if err := json.Unmarshal(content, &user); err != nil {
				return
			}
			matched := strings.Contains(user.Email, "."+domain)
			if matched {
				domain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
				cnt.mux.Lock()
				cnt.v[domain]++
				cnt.mux.Unlock()
			}
		}(append([]byte(nil), scanner.Bytes()...))
	}
	wg.Wait()
	return cnt.v, nil
}
