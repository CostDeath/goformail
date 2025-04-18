package db

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"log"
)

type Error struct {
	Err  error
	Code ErrorCode
}

type ErrorCode int

const (
	Unknown ErrorCode = iota
	ErrDuplicate
	ErrNoRows
)

func getError(err error) *Error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			return &Error{err, ErrDuplicate}
		default:
			log.Println(err)
			return &Error{err, Unknown}
		}
	} else if errors.Is(err, sql.ErrNoRows) {
		return &Error{err, ErrNoRows}
	}

	log.Println(err)
	return &Error{err, Unknown}
}
