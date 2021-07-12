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
	pos := strings.Index(line, "\"Email\":")
	if pos < 0 {
		return "", fmt.Errorf("bad input json: %s", line)
	}
	line = line[pos+8:]

	parts := strings.SplitN(line, "\"", 3)
	if len(parts) < 3 {
		return "", fmt.Errorf("bad input json: %s", line)
	}

	return parts[1], nil
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
