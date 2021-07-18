package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	tagLen    = "len"
	tagRegExp = "regexp"
	tagIn     = "in"
	tagMin    = "min"
	tagMax    = "max"
)

type checkerFunc func(value reflect.Value, fieldName, rule string) ValidationErrors

var stringValidator = map[string]checkerFunc{
	tagLen:    checkLen,
	tagRegExp: checkRegExp,
	tagIn:     checkIn,
}

var intValidator = map[string]checkerFunc{
	tagMin: checkMin,
	tagMax: checkMax,
	tagIn:  checkIn,
}

func checkLen(value reflect.Value, fieldName, rule string) ValidationErrors {
	l, err := strconv.Atoi(rule)
	if err != nil {
		return newSingleValidationErrors(
			fieldName,
			fmt.Errorf("%w: length value must be number", ErrTagParse),
		)
	}
	if len(value.String()) > l {
		return newSingleValidationErrors(
			fieldName,
			fmt.Errorf("%w: length is greater than %d", ErrValidation, l),
		)
	}
	return nil
}

func checkRegExp(value reflect.Value, fieldName, rule string) ValidationErrors {
	if len(rule) == 0 {
		return newSingleValidationErrors(
			fieldName,
			fmt.Errorf("%w: empty regexp", ErrTagParse),
		)
	}
	re, err := regexp.Compile(rule)
	if err != nil {
		return newSingleValidationErrors(
			fieldName,
			fmt.Errorf("%w: bad regexp %s", ErrTagParse, rule),
		)
	}
	if !re.MatchString(value.String()) {
		return newSingleValidationErrors(
			fieldName,
			fmt.Errorf("%w: not matched by regexp %s", ErrValidation, rule),
		)
	}
	return nil
}

func checkMin(value reflect.Value, fieldName, rule string) ValidationErrors {
	min, err := strconv.Atoi(rule)
	if err != nil {
		return newSingleValidationErrors(
			fieldName,
			fmt.Errorf("%w: min value must be number", ErrTagParse),
		)
	}

	if value.Int() < int64(min) {
		return newSingleValidationErrors(
			fieldName,
			fmt.Errorf("%w: must greater or equals than %d", ErrValidation, min),
		)
	}
	return nil
}

func checkMax(value reflect.Value, fieldName, rule string) ValidationErrors {
	max, err := strconv.Atoi(rule)
	if err != nil {
		return newSingleValidationErrors(
			fieldName,
			fmt.Errorf("%w: max value must be number", ErrTagParse),
		)
	}

	if value.Int() > int64(max) {
		return newSingleValidationErrors(
			fieldName,
			fmt.Errorf("%w: must lower or equals than %d", ErrValidation, max),
		)
	}
	return nil
}

func checkIn(value reflect.Value, fieldName, rule string) ValidationErrors {
	rule = strings.TrimSpace(rule)
	if len(rule) == 0 {
		return newSingleValidationErrors(
			fieldName,
			fmt.Errorf("%w: empty values for validation", ErrTagParse),
		)
	}
	var searchValue string
	switch value.Kind() {
	case reflect.Int:
		searchValue = strconv.FormatInt(value.Int(), 10)
	case reflect.String:
		searchValue = value.String()
	default:
		return newSingleValidationErrors(
			fieldName,
			fmt.Errorf("%w: unsupported type - %s", ErrTagParse, value.Kind()),
		)
	}
	if len(searchValue) == 0 {
		return newSingleValidationErrors(
			fieldName,
			fmt.Errorf("%w: empty value", ErrValidation),
		)
	}
	inValues := strings.Split(rule, ",")
	for _, s := range inValues {
		if searchValue == strings.TrimSpace(s) {
			return nil
		}
	}
	return newSingleValidationErrors(
		fieldName,
		fmt.Errorf("%w: unexpected value - %s", ErrValidation, searchValue),
	)
}
