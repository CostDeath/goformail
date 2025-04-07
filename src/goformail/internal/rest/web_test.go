package rest

//
//import (
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func cleanUp() {
//	http.DefaultServeMux = new(http.ServeMux)
//}
//
//func TestAddUiHandlerAddsUIEndpoint(t *testing.T) {
//	if err := addWebUiHandler(); err != nil {
//		t.Fatalf("Error adding web ui: %v", err)
//	}
//
//	// Mock the request
//	req, err := http.NewRequest("GET", "/ui/", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	rr := httptest.NewRecorder()
//	http.DefaultServeMux.ServeHTTP(rr, req)
//
//	// Check the status code is what we expect.
//	if rr.Code != http.StatusOK {
//		t.Errorf("Handler returned wrong status code: got %v wanted %v",
//			rr.Code, http.StatusOK)
//	}
//
//	// Check the response body is what we expect.
//	expected := "<!doctype html>\n<meta name=\"viewport\" content=\"width=device-width\">\n<pre>\n<a href=\".gitkeep\">.gitkeep</a>\n</pre>\n"
//	if rr.Body.String() != expected {
//		t.Errorf("Handler returned unexpected body: got %v wanted %v",
//			rr.Body.String(), expected)
//	}
//	t.Cleanup(cleanUp)
//}
//
//func TestAddUiHandlerAddsRootRedirect(t *testing.T) {
//	if err := addWebUiHandler(); err != nil {
//		t.Fatalf("Error adding web ui: '%v'", err)
//	}
//
//	// Mock the request
//	req, err := http.NewRequest("GET", "/", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	rr := httptest.NewRecorder()
//	http.DefaultServeMux.ServeHTTP(rr, req)
//
//	// Check the status code is what we expect.
//	if rr.Code != http.StatusFound {
//		t.Errorf("Handler returned wrong status code: got '%v' wanted '%v'",
//			rr.Code, http.StatusFound)
//	}
//
//	// Check the response body is what we expect.
//	expected := "<a href=\"/ui/\">Found</a>.\n\n"
//	if rr.Body.String() != expected {
//		t.Errorf("Handler returned unexpected body: got '%v' wanted '%v'",
//			rr.Body.String(), expected)
//	}
//	t.Cleanup(cleanUp)
//}
//
//func TestAddUiHandlerAdds404Page(t *testing.T) {
//	if err := addWebUiHandler(); err != nil {
//		t.Fatalf("Error adding web ui: '%v'", err)
//	}
//
//	// Mock the request
//	req, err := http.NewRequest("GET", "/invalidEndpoint", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	rr := httptest.NewRecorder()
//	http.DefaultServeMux.ServeHTTP(rr, req)
//
//	// Check the status code is what we expect.
//	if rr.Code != http.StatusNotFound {
//		t.Errorf("Handler returned wrong status code: got '%v' wanted '%v'",
//			rr.Code, http.StatusNotFound)
//	}
//
//	// Check the response body is what we expect.
//	expected := "404 page not found\n"
//	if rr.Body.String() != expected {
//		t.Errorf("Handler returned unexpected body: got '%v' wanted '%v'",
//			rr.Body.String(), expected)
//	}
//	t.Cleanup(cleanUp)
//}
