package rest

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetsResponseWithStringData(t *testing.T) {
	rr := httptest.NewRecorder()
	setResponse("test message", "data", rr, nil)

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	actualBody := rr.Body.String()
	expectedBody, err := json.Marshal(response{Message: "test message", Data: "data"})
	require.NoError(t, err)
	assert.Equal(t, string(expectedBody), actualBody)
}

func TestSetsResponseWithStructData(t *testing.T) {
	rr := httptest.NewRecorder()
	setResponse("test message", IdObject{Id: 1}, rr, nil)

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	actualBody := rr.Body.String()
	expectedBody, err := json.Marshal(response{Message: "test message", Data: IdObject{Id: 1}})
	require.NoError(t, err)
	assert.Equal(t, string(expectedBody), actualBody)
}

func TestHandleUnknownMethod(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/invalid", nil)

	handleUnknownMethod(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "No known method GET for endpoint /api/invalid\n", rr.Body.String())
}
