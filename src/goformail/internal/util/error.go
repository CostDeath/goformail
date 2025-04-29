package util

import (
	"fmt"
)

type Error struct {
	Err     error
	Code    ErrorCode
	Message string
}

type ErrorCode int

const (
	Unknown ErrorCode = iota
	ErrInvalidObject
	ErrNoUser
	ErrUserAlreadyExists
)

func NewGenericError(err error) *Error {
	msg := "Unknown error occurred"
	if err != nil {
		msg = fmt.Sprintf("Unknown error occurred: %s", err.Error())
	}
	return &Error{Err: err, Code: Unknown, Message: msg}
}

func NewInvalidObjectError(reason string, err error) *Error {
	msg := fmt.Sprintf("Invalid object: %s", reason)
	return &Error{Err: err, Code: ErrInvalidObject, Message: msg}
}

func NewNoUserError(id int, err error) *Error {
	msg := fmt.Sprintf("Could not find a user with id '%d'", id)
	return &Error{Err: err, Code: ErrNoUser, Message: msg}
}

func NewUserAlreadyExistsError(email string, err error) *Error {
	msg := fmt.Sprintf("A user with the email '%s' already exists", email)
	return &Error{Err: err, Code: ErrUserAlreadyExists, Message: msg}
}
