package util

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestSetsResponseWithStringData(t *testing.T) {
	rr := httptest.NewRecorder()
	SetResponse("test message", "data", rr, nil)

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	actualBody := rr.Body.String()
	expectedBody, err := json.Marshal(response{Message: "test message", Data: "data"})
	if err != nil || actualBody != string(expectedBody) {
		t.Errorf("Incorrect response body. Expected: '%s', got '%s'", expectedBody, actualBody)
	}
}

func TestSetsResponseWithStructData(t *testing.T) {
	rr := httptest.NewRecorder()
	SetResponse("test message", IdObject{Id: 1}, rr, nil)

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	actualBody := rr.Body.String()
	expectedBody, err := json.Marshal(response{Message: "test message", Data: IdObject{Id: 1}})
	if err != nil || actualBody != string(expectedBody) {
		t.Errorf("Incorrect response body. Expected: '%s', got '%s'", expectedBody, actualBody)
	}
}
