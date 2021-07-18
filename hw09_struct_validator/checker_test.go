package hw09structvalidator

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChecker_ValidRule_ValidValue(t *testing.T) {
	tests := []struct {
		checker checkerFunc
		input   interface{}
		rule    string
	}{
		{
			checker: checkLen,
			input:   "1234567890",
			rule:    "10",
		},
		{
			checker: checkRegExp,
			input:   "12345678901",
			rule:    "\\d+$",
		},
		{
			checker: checkMin,
			input:   15,
			rule:    "15",
		},
		{
			checker: checkMax,
			input:   20,
			rule:    "20",
		},
		{
			checker: checkIn,
			input:   20,
			rule:    "10, 20, 30",
		},
		{
			checker: checkIn,
			input:   "20",
			rule:    "10, 20, 30",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("case %s", runtime.FuncForPC(reflect.ValueOf(tc.checker).Pointer()).Name()), func(t *testing.T) {
			err := tc.checker(reflect.ValueOf(tc.input), "somefield", tc.rule)
			require.Nil(t, err)
		})
	}
}

func TestChecker_ValidRule_InvalidValue(t *testing.T) {
	tests := []struct {
		checker checkerFunc
		input   interface{}
		rule    string
	}{
		{
			checker: checkLen,
			input:   "12345678901",
			rule:    "10",
		},
		{
			checker: checkRegExp,
			input:   "12345678901d",
			rule:    "\\d+$",
		},
		{
			checker: checkMin,
			input:   14,
			rule:    "15",
		},
		{
			checker: checkMax,
			input:   21,
			rule:    "20",
		},
		{
			checker: checkIn,
			input:   21,
			rule:    "10, 20, 30",
		},
		{
			checker: checkIn,
			input:   "21",
			rule:    "10, 20, 30,,,",
		},
		{
			checker: checkIn,
			input:   "",
			rule:    "1,2",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("case %s", runtime.FuncForPC(reflect.ValueOf(tc.checker).Pointer()).Name()), func(t *testing.T) {
			fieldName := "somefield"
			err := tc.checker(reflect.ValueOf(tc.input), fieldName, tc.rule)
			assertFirstValidationError(t, err, fieldName, ErrValidation)
		})
	}
}

func TestChecker_InvalidRule_Error(t *testing.T) {
	tests := []struct {
		checker checkerFunc
		input   interface{}
		rule    string
	}{
		{
			checker: checkLen,
			input:   "12345678901",
			rule:    "",
		},
		{
			checker: checkLen,
			input:   "12345678901",
			rule:    "b",
		},
		{
			checker: checkLen,
			input:   "12345678901",
			rule:    "10b",
		},
		{
			checker: checkRegExp,
			input:   "12345678901d",
			rule:    "",
		},
		{
			checker: checkRegExp,
			input:   "12345678901d",
			rule:    "[[",
		},
		{
			checker: checkMin,
			input:   14,
			rule:    "",
		},
		{
			checker: checkMin,
			input:   14,
			rule:    "b",
		},
		{
			checker: checkMin,
			input:   14,
			rule:    "15b",
		},
		{
			checker: checkMax,
			input:   21,
			rule:    "",
		},
		{
			checker: checkMax,
			input:   21,
			rule:    "d",
		},
		{
			checker: checkMax,
			input:   21,
			rule:    "20s",
		},
		{
			checker: checkIn,
			input:   20,
			rule:    "",
		},
		{
			checker: checkIn,
			input:   true,
			rule:    "1,2",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("case %s", runtime.FuncForPC(reflect.ValueOf(tc.checker).Pointer()).Name()), func(t *testing.T) {
			fieldName := "somefield"
			err := tc.checker(reflect.ValueOf(tc.input), fieldName, tc.rule)
			assertFirstValidationError(t, err, fieldName, ErrTagParse)
		})
	}
}
