package util

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

func TestSetsResponseWithStringData(t *testing.T) {
	rr := httptest.NewRecorder()
	SetResponse("test message", "data", rr, nil)

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	actualBody := rr.Body.String()
	expectedBody, err := json.Marshal(response{Message: "test message", Data: "data"})
	require.NoError(t, err)
	assert.Equal(t, string(expectedBody), actualBody)
}

func TestSetsResponseWithStructData(t *testing.T) {
	rr := httptest.NewRecorder()
	SetResponse("test message", IdObject{Id: 1}, rr, nil)

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	actualBody := rr.Body.String()
	expectedBody, err := json.Marshal(response{Message: "test message", Data: IdObject{Id: 1}})
	require.NoError(t, err)
	assert.Equal(t, string(expectedBody), actualBody)
}
