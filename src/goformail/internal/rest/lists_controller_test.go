package rest

import (
	"bytes"
	"encoding/json"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/rest/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListPostEndpoint(t *testing.T) {
	mux := new(http.ServeMux)
	AddListHandlers(mux, test.Configs)

	// Mock the request
	list, err := json.Marshal(model.List{Name: "name"})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/api/list/", bytes.NewBuffer(list))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if rr.Code != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v wanted %v",
			rr.Code, http.StatusCreated)
	}

	// Check the response body is what we expect.
	expected := string(list)
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v wanted %v",
			rr.Body.String(), expected)
	}
}
