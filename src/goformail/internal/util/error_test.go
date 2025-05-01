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

func TestNewNoListError(t *testing.T) {
	err := errors.New("some error")
	expected := &Error{Err: err, Code: ErrNoList, Message: "Could not find a list with id '0'"}
	actual := NewNoListError(0, err)
	assert.Equal(t, expected, actual)
}

func TestNewListAlreadyExistsError(t *testing.T) {
	err := errors.New("some error")
	expected := &Error{Err: err, Code: ErrListAlreadyExists,
		Message: "A list with the name 'test' already exists"}
	actual := NewListAlreadyExistsError("test", err)
	assert.Equal(t, expected, actual)
}

func TestNewNoUserError(t *testing.T) {
	err := errors.New("some error")
	expected := &Error{Err: err, Code: ErrNoUser, Message: "Could not find a user with id '0'"}
	actual := NewNoUserError(0, err)
	assert.Equal(t, expected, actual)
}

func TestNewNoUserEmailError(t *testing.T) {
	err := errors.New("some error")
	expected := &Error{Err: err, Code: ErrNoUser, Message: "Could not find a user with email 'test@domain.tld'"}
	actual := NewNoUserEmailError("test@domain.tld", err)
	assert.Equal(t, expected, actual)
}

func TestNewUserAlreadyExistsError(t *testing.T) {
	err := errors.New("some error")
	expected := &Error{Err: err, Code: ErrUserAlreadyExists,
		Message: "A user with the email 'test@domain.tld' already exists"}
	actual := NewUserAlreadyExistsError("test@domain.tld", err)
	assert.Equal(t, expected, actual)
}

func TestNewIncorrectPasswordError(t *testing.T) {
	err := errors.New("some error")
	expected := &Error{Err: err, Code: ErrIncorrectPassword,
		Message: "Incorrect password for user 'test@domain.tld'"}
	actual := NewIncorrectPasswordError("test@domain.tld", err)
	assert.Equal(t, expected, actual)
}

func TestNewInvalidTokenError(t *testing.T) {
	err := errors.New("some error")
	expected := &Error{Err: err, Code: ErrInvalidToken,
		Message: "Invalid token provided"}
	actual := NewInvalidTokenError(err)
	assert.Equal(t, expected, actual)
}

func TestNewNoPermissionError(t *testing.T) {
	err := errors.New("some error")
	expected := &Error{Err: err, Code: ErrNoPermission,
		Message: "Missing permission 'CRT_LIST' for this action"}
	actual := NewNoPermissionError("CRT_LIST", err)
	assert.Equal(t, expected, actual)
}

func TestNewEncryptionError(t *testing.T) {
	err := errors.New("some error")
	expected := &Error{Err: err, Code: ErrEncryption,
		Message: "An error relating to encryption occurred: 'some error'"}
	actual := NewEncryptionError(err)
	assert.Equal(t, expected, actual)
}
