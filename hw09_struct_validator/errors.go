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
	errStr := ""
	for _, err := range v {
		errStr = fmt.Sprintf("%s%s: %s, ", errStr, err.Field, err.Err)
	}

	return strings.TrimSuffix(errStr, ", ")
}

func (v ValidationError) Error() string {
	return v.Err.Error()
}
