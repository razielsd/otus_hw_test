package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:20"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	BadTagUser struct {
		ID string `json:"id" validate:"len10"`
	}
)

func TestValidate_ValidValues(t *testing.T) {
	tests := []struct {
		in interface{}
	}{
		{
			in: User{
				ID:     "12345-67890-12-12345",
				Name:   "Petr",
				Age:    50,
				Email:  "me@example.com",
				Role:   "admin",
				Phones: []string{"+7916123456", "+7916123457"},
			},
		},
		{
			in: App{
				Version: "1.0.1",
			},
		},
		{
			in: Token{
				Header:    []byte("header here"),
				Payload:   []byte("{}"),
				Signature: []byte("12sda121211289"),
			},
		},
		{
			in: Response{
				Code: 200,
			},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tc := tc
			err := Validate(tc.in)
			require.NoError(t, err)
		})
	}
}

func TestValidate_BadValues(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr ValidationErrors
	}{
		{
			in: User{
				ID:     "12345-67890-12-123451",
				Name:   "Petr",
				Age:    51,
				Email:  "me@",
				Role:   "adminchik",
				Phones: []string{"+791612345611", "+7916123457"},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				newValidationError("ID", ErrValidation),
				newValidationError("Age", ErrValidation),
				newValidationError("Email", ErrValidation),
				newValidationError("Role", ErrValidation),
				newValidationError("Phones[0]", ErrValidation),
			},
		},
		{
			in: App{
				Version: "1.0.1bv",
			},
			expectedErr: ValidationErrors{
				newValidationError("Version", ErrValidation),
			},
		},
		{
			in: Response{
				Code: 502,
			},
			expectedErr: ValidationErrors{
				newValidationError("Code", ErrValidation),
			},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tc := tc
			err := Validate(tc.in)
			var ve ValidationErrors
			require.ErrorAs(t, err, &ve)

			extractName := func(v ValidationErrors) []string {
				var r []string
				for _, f := range v {
					r = append(r, f.Field)
				}
				return r
			}
			require.Equal(t, extractName(tc.expectedErr), extractName(ve), "Not all errors found")
			for i, fErr := range ve {
				exp := tc.expectedErr[i]
				require.Equal(t, exp.Field, fErr.Field)
				require.ErrorIs(t, fErr.Err, exp.Err)
			}
		})
	}
}

func TestValidate_NotStruct_Error(t *testing.T) {
	err := Validate(1)
	require.ErrorIs(t, err, ErrRequireStruct)
}

func TestIsUpperFirst_IsUpper(t *testing.T) {
	tests := []string{"U", "Upper", "Б", "Большой"}
	for i := range tests {
		t.Run(fmt.Sprintf("case %s", tests[i]), func(t *testing.T) {
			require.True(t, isUpperFirst(tests[i]))
		})
	}
}

func TestIsUpperFirst_NotUpper(t *testing.T) {
	tests := []string{"", "l", "lower", "м", "маленький", "lOweR", "мАленькиЙ", "_", "1"}
	for i := range tests {
		t.Run(fmt.Sprintf("case %s", tests[i]), func(t *testing.T) {
			require.False(t, isUpperFirst(tests[i]))
		})
	}
}

func TestParseTag_ValidTag_Success(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]string
	}{
		{
			input:    "",
			expected: map[string]string{},
		},
		{
			input:    "len:10",
			expected: map[string]string{"len": "10"},
		},
		{
			input:    "Len:10",
			expected: map[string]string{"len": "10"},
		},
		{
			input:    "min:10|max:40",
			expected: map[string]string{"min": "10", "max": "40"},
		},
		{
			input:    "in:admin,stuff",
			expected: map[string]string{"in": "admin,stuff"},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("case %s", tc.input), func(t *testing.T) {
			actual, err := parseTag(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestParseTag_InvalidTag_Error(t *testing.T) {
	tests := []struct {
		input string
	}{
		{input: "len10"},
		{input: "min10|max:40"},
		{input: "inadmin,stuff"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("case %s", tc.input), func(t *testing.T) {
			_, err := parseTag(tc.input)
			require.ErrorIs(t, err, ErrTagParse)
		})
	}
}

func TestValidateBasicTypeField_UnknownType(t *testing.T) {
	var v byte
	tagInfo := map[string]string{"len": "10"}
	err := validateBasicTypeField("Debug", reflect.ValueOf(v), tagInfo)
	require.Nil(t, err)
}

func TestValidateBasicTypeField_UnknownChecker(t *testing.T) {
	tagInfo := map[string]string{"debug": "123"}
	err := validateBasicTypeField(
		"Debug", reflect.ValueOf(55), tagInfo,
	)
	assertFirstValidationError(t, err, "Debug", ErrTagParse)
}

func TestValidate_BadCheckInTag(t *testing.T) {
	bad := BadTagUser{
		ID: "1234",
	}
	err := Validate(bad)
	assertFirstValidationError(t, err, "ID", ErrTagParse)
}

func assertFirstValidationError(t *testing.T, err error, fieldName string, expErr error) {
	var ve ValidationErrors
	require.ErrorAs(t, err, &ve)
	require.Len(t, ve, 1)
	fieldErr := ve[0]
	require.Error(t, fieldErr, "error not found for field")
	require.Equal(t, fieldName, fieldErr.Field)
	require.ErrorIs(t, fieldErr.Err, expErr)
}
