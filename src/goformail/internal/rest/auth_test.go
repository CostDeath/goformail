package rest

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

var defaultLoginRequest = &model.LoginRequest{Email: "example@domain.tld", Password: "pass"}

func authCleanUp() {
	http.DefaultServeMux = new(http.ServeMux)
}

func TestLogin(t *testing.T) {
	t.Cleanup(authCleanUp)
	mockObj := new(service.IAuthManagerMock)
	mockObj.On("Login", defaultLoginRequest).Return("token", 1)
	ctrl := &Controller{auth: mockObj, mux: new(http.ServeMux)}
	ctrl.addAuthHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/login/", defaultLoginRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	mockObj.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := test.GetExpectedJsonResponse(t, "Successfully logged in!", LoginObject{Token: "token", User: 1})
	assert.Equal(t, expected, rr.Body.String())
}

func TestLogin400sOnInvalidJson(t *testing.T) {
	t.Cleanup(authCleanUp)
	mockObj := new(service.IAuthManagerMock)
	mockObj.On("Login", mock.Anything, mock.Anything).Panic("Login should not have been called")
	ctrl := &Controller{auth: mockObj, mux: new(http.ServeMux)}
	ctrl.addAuthHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/login/", 1)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid json provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestLogin404sOnNoUser(t *testing.T) {
	t.Cleanup(authCleanUp)
	mockObj := service.NewIAuthManagerMockWithError(util.ErrNoUser)
	mockObj.On("Login", defaultLoginRequest).Return("", 0)
	ctrl := &Controller{auth: mockObj, mux: new(http.ServeMux)}
	ctrl.addAuthHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/login/", defaultLoginRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestLogin401sOnIncorrectPassword(t *testing.T) {
	t.Cleanup(authCleanUp)
	mockObj := service.NewIAuthManagerMockWithError(util.ErrIncorrectPassword)
	mockObj.On("Login", defaultLoginRequest).Return("", 0)
	ctrl := &Controller{auth: mockObj, mux: new(http.ServeMux)}
	ctrl.addAuthHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/login/", defaultLoginRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestLogin500sOnGenericError(t *testing.T) {
	t.Cleanup(authCleanUp)
	mockObj := service.NewIAuthManagerMockWithError(util.Unknown)
	mockObj.On("Login", defaultLoginRequest).Return("", 0)
	ctrl := &Controller{auth: mockObj, mux: new(http.ServeMux)}
	ctrl.addAuthHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/login/", defaultLoginRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestValidate(t *testing.T) {
	t.Cleanup(authCleanUp)
	mockObj := new(service.IAuthManagerMock)
	mockObj.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{auth: mockObj, mux: new(http.ServeMux)}
	ctrl.addAuthHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/validateToken/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	mockObj.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := test.GetExpectedJsonResponse(t, "Token is valid!", TokenObject{Token: "token"})
	assert.Equal(t, expected, rr.Body.String())
}

func TestValidate401sOnInvalidToken(t *testing.T) {
	t.Cleanup(authCleanUp)
	mockObj := service.NewIAuthManagerMockWithError(util.ErrInvalidToken)
	mockObj.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{auth: mockObj, mux: new(http.ServeMux)}
	ctrl.addAuthHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/validateToken/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestValidate500sOnGenericError(t *testing.T) {
	t.Cleanup(authCleanUp)
	mockObj := service.NewIAuthManagerMockWithError(util.Unknown)
	mockObj.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{auth: mockObj, mux: new(http.ServeMux)}
	ctrl.addAuthHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/validateToken/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}
