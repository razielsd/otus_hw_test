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
		email, err := ExtractEmail(line)
		if err != nil {
			return nil, fmt.Errorf("bad input json: %s", line)
		}
		host, ok := extractor(email)
		if !ok {
			continue
		}
		result[host]++
	}
	return result, nil
}

func ExtractEmail(line string) (string, error) {
	c := strings.Split(line, "\"Email\":")
	if len(c) < 2 {
		return "", fmt.Errorf("bad input json: %s", line)
	}
	c = strings.Split(c[1], "\",")
	return c[0], nil
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
