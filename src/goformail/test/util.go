package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func CreateHttpRequest(t *testing.T, method string, uri string, body interface{}) *http.Request {
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	return req
}

func GetExpectedJsonResponse(t *testing.T, msg string, data interface{}) string {
	jsonData, err := json.Marshal(data)
	require.NoError(t, err)
	return fmt.Sprintf("{\"message\":\"%s\",\"data\":%s}", msg, jsonData)
}
