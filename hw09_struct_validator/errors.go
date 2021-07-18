package hw09structvalidator

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func newValidationErrors() ValidationErrors {
	return ValidationErrors{}
}

func newSingleValidationErrors(name string, err error) ValidationErrors {
	return ValidationErrors{
		newValidationError(name, err),
	}
}

func newValidationError(fieldName string, err error) ValidationError {
	return ValidationError{
		Field: fieldName,
		Err:   err,
	}
}

func (v ValidationErrors) Error() string {
	b := strings.Builder{}
	for _, err := range v {
		b.WriteString(fmt.Sprintf("%s: %s, ", err.Field, err.Err))
	}

	return strings.TrimSuffix(b.String(), ", ")
}

func (v ValidationError) Error() string {
	return v.Err.Error()
}
