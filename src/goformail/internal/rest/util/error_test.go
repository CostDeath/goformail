package util

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleErrorInvalidMethod(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/invalid", nil)
	HandleError(rr, req, InvalidMethod)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "No known method GET for endpoint /invalid\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestHandleErrorInvalidPayload(t *testing.T) {
	rr := httptest.NewRecorder()
	HandleError(rr, nil, InvalidPayload)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid payload\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestHandleErrorListAlreadyExists(t *testing.T) {
	rr := httptest.NewRecorder()
	HandleError(rr, nil, ListAlreadyExists)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusConflict, rr.Code)
	expected := "A list with this name already exists\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestHandleErrorListNotFound(t *testing.T) {
	rr := httptest.NewRecorder()
	HandleError(rr, nil, ListNotFound)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "List not found\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestHandleErrorDatabaseError(t *testing.T) {
	rr := httptest.NewRecorder()
	HandleError(rr, nil, DatabaseError)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "An error occurred relating to the database\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestHandleErrorGeneric(t *testing.T) {
	rr := httptest.NewRecorder()
	HandleError(rr, nil, 9999)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "An unknown error occurred\n"
	assert.Equal(t, expected, rr.Body.String())
}
