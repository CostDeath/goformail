package db

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetErrorReturnsDuplicate(t *testing.T) {
	pqErr := &pq.Error{Code: "23505"}
	actual := getError(pqErr)
	assert.Equal(t, &Error{Err: pqErr, Code: ErrDuplicate}, actual)
}

func TestGetErrorReturnsNoRows(t *testing.T) {
	sqlErr := sql.ErrNoRows
	actual := getError(sqlErr)
	assert.Equal(t, &Error{Err: sqlErr, Code: ErrNoRows}, actual)
}

func TestGetErrorReturnsUnknownWhenUnknownCode(t *testing.T) {
	pqErr := &pq.Error{Code: "0", Message: "unknown pq error"}
	actual := getError(pqErr)
	assert.Equal(t, &Error{Err: pqErr, Code: Unknown}, actual)
}

func TestGetErrorReturnsUnknownWhenUnknownErr(t *testing.T) {
	err := errors.New("unknown error")
	actual := getError(err)
	assert.Equal(t, &Error{Err: err, Code: Unknown}, actual)
}
