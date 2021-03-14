package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode/utf8"
)

const (
	CharTypeUndefined = iota
	CharTypeSymbol
	CharTypeEscape
	CharTypeNumber
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(encodedString string) (string, error) {
	if len(encodedString) == 0 {
		return "", nil
	}
	if !utf8.ValidString(encodedString) {
		return "", ErrInvalidString
	}

	var b strings.Builder
	var prevCharValue rune
	prevCharType := CharTypeUndefined
	encodedString += " " // last symbol ignored, add one new
	for i, char := range encodedString {
		charType, number := parseRune(char)
		if i == 0 && charType == CharTypeNumber {
			return "", ErrInvalidString
		}
		if charType == CharTypeNumber && prevCharType == CharTypeNumber {
			return "", ErrInvalidString
		}
		if prevCharType == CharTypeEscape {
			charType = CharTypeSymbol
		}
		if prevCharType == CharTypeSymbol {
			b.WriteString(strings.Repeat(string(prevCharValue), number))
		}

		prevCharValue = char
		prevCharType = charType
	}

	return b.String(), nil
}

func parseRune(char rune) (int, int) {
	if char == '\\' {
		return CharTypeEscape, 1
	} else if char < '0' || char > '9' {
		return CharTypeSymbol, 1
	}

	number := int(char - '0')
	return CharTypeNumber, number
}
