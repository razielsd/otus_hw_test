package hw09structvalidator

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSingleValidationErrors(t *testing.T) {
	fieldName := "myfield"
	fieldErr := errors.New("my field error")
	err := newSingleValidationErrors(fieldName, fieldErr)
	require.ErrorAs(t, err, &ValidationErrors{})
	require.Len(t, err, 1)
	require.Equal(t, fieldName, err[0].Field)
	require.Equal(t, fieldErr, err[0].Err)
}

func TestValidationError_Error(t *testing.T) {
	fieldName := "myField"
	fieldErr := errors.New("my field error")
	err := ValidationError{
		Field: fieldName,
		Err:   fieldErr,
	}
	require.Equal(
		t,
		fieldErr.Error(),
		err.Error(),
	)
}

func TestValidationErrors_Error(t *testing.T) {
	errs := ValidationErrors{}

	expTxt := ""
	for i := 0; i < 3; i++ {
		fieldName := fmt.Sprintf("myField%d", i)
		fieldErr := fmt.Errorf("my field error %d", i)
		err := ValidationError{
			Field: fieldName,
			Err:   fieldErr,
		}
		errs = append(errs, err)
		expTxt += fmt.Sprintf("%s: %s, ", fieldName, fieldErr)
	}
	expTxt = strings.TrimSuffix(expTxt, ", ")
	require.Equal(t, expTxt, errs.Error())
}
