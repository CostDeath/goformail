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
	ErrNoList
	ErrListAlreadyExists
	ErrNoUser
	ErrUserAlreadyExists
	ErrNoEmail
	ErrIncorrectPassword
	ErrInvalidToken
	ErrNoPermission
	ErrEncryption
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

func NewNoListError(id int, err error) *Error {
	msg := fmt.Sprintf("Could not find a list with id '%d'", id)
	return &Error{Err: err, Code: ErrNoList, Message: msg}
}

func NewListAlreadyExistsError(name string, err error) *Error {
	msg := fmt.Sprintf("A list with the name '%s' already exists", name)
	return &Error{Err: err, Code: ErrListAlreadyExists, Message: msg}
}

func NewNoUserError(id int, err error) *Error {
	msg := fmt.Sprintf("Could not find a user with id '%d'", id)
	return &Error{Err: err, Code: ErrNoUser, Message: msg}
}

func NewNoUserEmailError(email string, err error) *Error {
	msg := fmt.Sprintf("Could not find a user with email '%s'", email)
	return &Error{Err: err, Code: ErrNoUser, Message: msg}
}

func NewUserAlreadyExistsError(email string, err error) *Error {
	msg := fmt.Sprintf("A user with the email '%s' already exists", email)
	return &Error{Err: err, Code: ErrUserAlreadyExists, Message: msg}
}

func NewNoEmailError(id int, err error) *Error {
	msg := fmt.Sprintf("Could not find a email with id '%d'", id)
	return &Error{Err: err, Code: ErrNoEmail, Message: msg}
}

func NewIncorrectPasswordError(email string, err error) *Error {
	msg := fmt.Sprintf("Incorrect password for user '%s'", email)
	return &Error{Err: err, Code: ErrIncorrectPassword, Message: msg}
}

func NewInvalidTokenError(err error) *Error {
	msg := "Invalid token provided"
	return &Error{Err: err, Code: ErrInvalidToken, Message: msg}
}

func NewNoPermissionError(perm string, err error) *Error {
	msg := fmt.Sprintf("Missing permission '%s' for this action", perm)
	return &Error{Err: err, Code: ErrNoPermission, Message: msg}
}

func NewEncryptionError(err error) *Error {
	msg := fmt.Sprintf("An error relating to encryption occurred: '%s'", err.Error())
	return &Error{Err: err, Code: ErrEncryption, Message: msg}
}
