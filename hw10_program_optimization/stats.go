package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	extractor := GetHostExtractor(domain)
	for scanner.Scan() {
		line := scanner.Text()
		c := strings.Split(line, "\"Email\":")
		if len(c) < 2 {
			return nil, fmt.Errorf("bad input json: %s", line)
		}
		c = strings.Split(c[1], "\",")
		host, ok := extractor(c[0])
		if !ok {
			continue
		}
		v := result[host]
		v++
		result[host] = v
	}
	return result, nil
}

func GetHostExtractor(domain string) func(email string) (string, bool) {
	sv := "." + domain
	return func(email string) (string, bool) {
		email = strings.ToLower(email)
		if strings.HasSuffix(email, sv) {
			info := strings.SplitN(email, "@", 2)
			if len(info) < 2 {
				return "", false
			}
			return info[1], true
		}
		return "", false
	}
}
