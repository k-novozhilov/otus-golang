package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

//easyjson:json
type User struct {
	ID       int    `json:"Id"`
	Name     string `json:"Name"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Phone    string `json:"Phone"`
	Password string `json:"Password"`
	Address  string `json:"Address"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	domainSuffix := "." + strings.ToLower(domain)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var user User
		if err := user.UnmarshalJSON(line); err != nil {
			return nil, fmt.Errorf("unmarshal error: %w", err)
		}

		emailParts := strings.SplitN(user.Email, "@", 2)
		if len(emailParts) == 2 {
			emailDomain := strings.ToLower(emailParts[1])
			if strings.HasSuffix(emailDomain, domainSuffix) {
				result[emailDomain]++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return result, nil
}
