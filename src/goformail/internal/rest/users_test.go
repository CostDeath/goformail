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

var defaultUserRequest = &model.UserRequest{
	Email: "example@domain.tld", Password: "pass", Permissions: []string{"ADMIN", "CRT_LIST"},
}

var defaultUserResponse = &model.UserResponse{
	Id: 1, Email: defaultUserRequest.Email, Permissions: defaultUserRequest.Permissions,
}

func usersCleanUp() {
	http.DefaultServeMux = new(http.ServeMux)
}

func TestGetUser(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := new(service.IUserManagerMock)
	mockObj.On("GetUser", 1).Return(defaultUserResponse)
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/user/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	mockObj.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := test.GetExpectedJsonResponse(t, "Successfully fetched user!", defaultUserResponse)
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetUser400sOnNoParam(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := new(service.IUserManagerMock)
	mockObj.On("GetUser", mock.Anything).Panic("GetUser should not have been called")
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/user/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid id provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetUser400sOnInvalidParam(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := new(service.IUserManagerMock)
	mockObj.On("GetUser", mock.Anything).Panic("GetUser should not have been called")
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/user/?id=a", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid id provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetUser404sWhenNoUser(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := service.NewIUserManagerMockWithError(util.ErrNoUser)
	mockObj.On("GetUser", 1).Return(defaultUserResponse)
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/user/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetUser500sOnGenericError(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := service.NewIUserManagerMockWithError(util.Unknown)
	mockObj.On("GetUser", 1).Return(defaultUserResponse)
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/user/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPostUser(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := new(service.IUserManagerMock)
	mockObj.On("CreateUser", defaultUserRequest).Return(1)
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/user/", defaultUserRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	mockObj.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusCreated, rr.Code)
	expected := test.GetExpectedJsonResponse(t, "Successfully created user!", IdObject{Id: 1})
	assert.Equal(t, expected, rr.Body.String())
}

func TestPostUser400sOnInvalidJson(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := new(service.IUserManagerMock)
	mockObj.On("CreateUser", mock.Anything, mock.Anything).Panic("UpdateUser should not have been called")
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/user/", 1)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid json provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPostUser400sOnInvalidObj(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := service.NewIUserManagerMockWithError(util.ErrInvalidObject)
	mockObj.On("CreateUser", defaultUserRequest).Return(1)
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/user/", defaultUserRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPostUser409sOnDuplicateUser(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := service.NewIUserManagerMockWithError(util.ErrUserAlreadyExists)
	mockObj.On("CreateUser", defaultUserRequest).Return(1)
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/user/", defaultUserRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusConflict, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPostUser500sOnDbError(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := service.NewIUserManagerMockWithError(util.Unknown)
	mockObj.On("CreateUser", defaultUserRequest).Return(1)
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/user/", defaultUserRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchUser(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := new(service.IUserManagerMock)
	mockObj.On("UpdateUser", 1, defaultUserRequest).Return(1)
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "PATCH", "/api/user/?id=1", defaultUserRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	mockObj.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := test.GetExpectedJsonResponse(t, "Successfully patched user!", IdObject{Id: 1})
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchUser400sOnNoParam(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := new(service.IUserManagerMock)
	mockObj.On("UpdateUser", mock.Anything, mock.Anything).Panic("UpdateUser should not have been called")
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "PATCH", "/api/user/", defaultUserRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid id provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchUser400sOnInvalidParam(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := new(service.IUserManagerMock)
	mockObj.On("UpdateUser", mock.Anything, mock.Anything).Panic("UpdateUser should not have been called")
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "PATCH", "/api/user/?id=a", defaultUserRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid id provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchUser400sOnInvalidJson(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := new(service.IUserManagerMock)
	mockObj.On("UpdateUser", mock.Anything, mock.Anything).Panic("UpdateUser should not have been called")
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "PATCH", "/api/user/?id=1", 1)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid json provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchUser400sOnInvalidObj(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := service.NewIUserManagerMockWithError(util.ErrInvalidObject)
	mockObj.On("UpdateUser", 1, defaultUserRequest).Return(1)
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "PATCH", "/api/user/?id=1", defaultUserRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchUser409sOnDuplicateName(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := service.NewIUserManagerMockWithError(util.ErrUserAlreadyExists)
	mockObj.On("UpdateUser", 1, defaultUserRequest).Return(1)
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "PATCH", "/api/user/?id=1", defaultUserRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusConflict, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchUser404sWhenNoUser(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := service.NewIUserManagerMockWithError(util.ErrNoUser)
	mockObj.On("UpdateUser", 1, defaultUserRequest).Return(1)
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "PATCH", "/api/user/?id=1", defaultUserRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchUser500sOnGenericError(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := service.NewIUserManagerMockWithError(util.Unknown)
	mockObj.On("UpdateUser", 1, defaultUserRequest).Return(1)
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "PATCH", "/api/user/?id=1", defaultUserRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestDeleteUser(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := new(service.IUserManagerMock)
	mockObj.On("DeleteUser", 1).Return()
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "DELETE", "/api/user/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	mockObj.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := test.GetExpectedJsonResponse(t, "Successfully deleted user!", IdObject{Id: 1})
	assert.Equal(t, expected, rr.Body.String())
}

func TestDeleteUser400sOnNoParam(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := new(service.IUserManagerMock)
	mockObj.On("DeleteUser", mock.Anything).Panic("DeleteUser should not have been called")
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "DELETE", "/api/user/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid id provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestDeleteUser404sOnInvalidParam(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := new(service.IUserManagerMock)
	mockObj.On("DeleteUser", mock.Anything).Panic("DeleteUser should not have been called")
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "DELETE", "/api/user/?id=a", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid id provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestDeleteUser404sWhenNoUser(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := service.NewIUserManagerMockWithError(util.ErrNoUser)
	mockObj.On("DeleteUser", 1).Return()
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "DELETE", "/api/user/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestDeleteUser500sOnGenericError(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := service.NewIUserManagerMockWithError(util.Unknown)
	mockObj.On("DeleteUser", 1).Return()
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "DELETE", "/api/user/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetUsers(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := new(service.IUserManagerMock)
	mockObj.On("GetAllUsers").Return(&[]*model.UserResponse{defaultUserResponse})
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/users/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	mockObj.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := test.GetExpectedJsonResponse(t, "Successfully fetched users!", []*model.UserResponse{defaultUserResponse})
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetUsers500sOnGenericError(t *testing.T) {
	t.Cleanup(usersCleanUp)
	mockObj := service.NewIUserManagerMockWithError(util.Unknown)
	mockObj.On("GetAllUsers").Return(&[]*model.UserResponse{defaultUserResponse})
	ctrl := NewController(util.MockConfigs, nil, mockObj)
	ctrl.addUserHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/users/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}
