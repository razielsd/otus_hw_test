package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Bytes()
		email, err := ExtractEmail(line)
		if err != nil {
			return nil, err
		}
		host, ok := ExtractHost(domain, email)
		if !ok {
			continue
		}
		result[host]++
	}
	return result, nil
}

func ExtractEmail(line []byte) (string, error) {
	if len(line) < 1 {
		return "", errors.New("empty json")
	}
	if err := fastjson.ValidateBytes(line); err != nil {
		return "", err
	}
	email := fastjson.GetString(line, "Email")
	// считаем, что у нас чаще есть поле Email, чем его нету
	if (email == "") && !fastjson.Exists(line, "Email") {
		return "", errors.New("not found file Email")
	}
	return email, nil
}

func ExtractHost(domain, email string) (string, bool) {
	email = strings.ToLower(email)
	if strings.HasSuffix(email, "."+domain) {
		info := strings.SplitN(email, "@", 2)
		if len(info) < 2 {
			return "", false
		}
		return info[1], true
	}
	return "", false
}
