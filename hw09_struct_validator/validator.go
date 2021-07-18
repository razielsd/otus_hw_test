package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

const (
	structKey = "validate"
)

var (
	ErrRequireStruct = errors.New("require struct")
	ErrTagParse      = errors.New("parse tag error")
	ErrValidation    = errors.New("validation error")
)

func Validate(v interface{}) error {
	refT := reflect.TypeOf(v)
	if refT.Kind() != reflect.Struct {
		return ErrRequireStruct
	}
	refV := reflect.ValueOf(v)

	ErrResult := ValidationErrors{}
	for i := 0; i < refT.NumField(); {
		field := refT.Field(i)
		i++
		if !isUpperFirst(field.Name) {
			continue
		}
		tag, ok := field.Tag.Lookup(structKey)
		if !ok {
			continue
		}
		tagInfo, pErr := parseTag(tag)
		if pErr != nil {
			ErrResult = append(ErrResult, newValidationError(field.Name, pErr))
			continue
		}
		err := validateField(field.Name, refV.FieldByName(field.Name), tagInfo)
		if err != nil {
			ErrResult = append(ErrResult, err...)
		}
	}
	if len(ErrResult) > 0 {
		return ErrResult
	}
	return nil
}

func validateField(fieldName string, value reflect.Value, tagInfo map[string]string) ValidationErrors {
	switch value.Kind() {
	case reflect.Slice:
		valErrors := newValidationErrors()
		for i := 0; i < value.Len(); i++ {
			value := value.Index(i)
			fieldName := fmt.Sprintf("%s[%d]", fieldName, i)
			err := validateBasicTypeField(fieldName, value, tagInfo)
			valErrors = append(valErrors, err...)
		}
		if len(valErrors) == 0 {
			return nil
		}
		return valErrors
	default:
		return validateBasicTypeField(fieldName, value, tagInfo)
	}
}

func validateBasicTypeField(fieldName string, value reflect.Value, tagInfo map[string]string) ValidationErrors {
	var TypeValidator map[string]checkerFunc
	for checkerName, v := range tagInfo {
		switch value.Kind() {
		case reflect.String:
			TypeValidator = stringValidator
		case reflect.Int:
			TypeValidator = intValidator
		default:
			return nil
		}
		validator, ok := TypeValidator[checkerName]
		if !ok {
			return newSingleValidationErrors(
				fieldName,
				fmt.Errorf("%w: unknown checker %s", ErrTagParse, checkerName),
			)
		}
		err := validator(value, fieldName, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func isUpperFirst(s string) bool {
	if len(s) == 0 {
		return false
	}
	return unicode.IsUpper([]rune(s)[0])
}

func parseTag(tag string) (map[string]string, error) {
	result := make(map[string]string)
	if len(tag) == 0 {
		return result, nil
	}

	parts := strings.Split(tag, "|")
	for _, v := range parts {
		vp := strings.SplitN(v, ":", 2)
		if len(vp) != 2 {
			return nil, ErrTagParse
		}
		result[strings.ToLower(vp[0])] = strings.TrimSpace(vp[1])
	}

	return result, nil
}
