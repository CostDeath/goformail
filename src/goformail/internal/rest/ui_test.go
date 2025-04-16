package rest

import (
	"github.com/stretchr/testify/assert"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"net/http"
	"net/http/httptest"
	"testing"
)

func uiCleanUp() {
	http.DefaultServeMux = new(http.ServeMux)
}

func TestAddUiHandlerAddsUIEndpoint(t *testing.T) {
	t.Cleanup(uiCleanUp)
	ctrl := NewController(util.MockConfigs, nil)
	ctrl.addUiHandler()

	// Mock the request
	req, err := http.NewRequest("GET", "/ui/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotEmpty(t, rr.Body.String())
}

func TestAddUiHandlerAddsRootRedirect(t *testing.T) {
	t.Cleanup(uiCleanUp)
	ctrl := NewController(util.MockConfigs, nil)
	ctrl.addUiHandler()

	// Mock the request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusFound, rr.Code)
	expected := "<a href=\"/ui/\">Found</a>.\n\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestAddUiHandlerAdds404Page(t *testing.T) {
	t.Cleanup(uiCleanUp)
	ctrl := NewController(util.MockConfigs, nil)
	ctrl.addUiHandler()

	// Mock the request
	req, err := http.NewRequest("GET", "/invalidEndpoint", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "404 page not found\n"
	assert.Equal(t, expected, rr.Body.String())
}
