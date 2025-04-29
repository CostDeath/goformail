package util

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGenericError(t *testing.T) {
	err := errors.New("some error")
	expected := &Error{Err: err, Code: Unknown, Message: "Unknown error occurred: some error"}
	actual := NewGenericError(err)
	assert.Equal(t, expected, actual)
}

func TestNewGenericErrorWithNilError(t *testing.T) {
	expected := &Error{Code: Unknown, Message: "Unknown error occurred"}
	actual := NewGenericError(nil)
	assert.Equal(t, expected, actual)
}

func TestNewInvalidObjectError(t *testing.T) {
	err := errors.New("some error")
	expected := &Error{Err: err, Code: ErrInvalidObject, Message: "Invalid object: some reason"}
	actual := NewInvalidObjectError("some reason", err)
	assert.Equal(t, expected, actual)
}

func TestNewNoUserError(t *testing.T) {
	err := errors.New("some error")
	expected := &Error{Err: err, Code: ErrNoUser, Message: "Could not find a user with id '0'"}
	actual := NewNoUserError(0, err)
	assert.Equal(t, expected, actual)
}

func TestNewUserAlreadyExistsError(t *testing.T) {
	err := errors.New("some error")
	expected := &Error{Err: err, Code: ErrUserAlreadyExists,
		Message: "A user with the email 'test@domain.tld' already exists"}
	actual := NewUserAlreadyExistsError("test@domain.tld", err)
	assert.Equal(t, expected, actual)
}
