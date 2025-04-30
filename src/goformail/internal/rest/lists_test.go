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

type listRequestWithoutLocked struct {
	Name            string   `json:"name"`
	Recipients      []string `json:"recipients"`
	Mods            []int64  `json:"mods"`
	ApprovedSenders []string `json:"approved_senders"`
}

var defaultListRequest = &model.ListRequest{
	Name:            "test",
	Recipients:      []string{"rcpt1@domain.tld", "rcpt2@domain.tld"},
	Mods:            []int64{1, 2},
	ApprovedSenders: []string{"sdr1@domain.tld", "sdr2@domain.tld"},
	Locked:          true,
}

var defaultListResponse = &model.ListResponse{
	Id:              1,
	Name:            defaultListRequest.Name,
	Recipients:      defaultListRequest.Recipients,
	Mods:            defaultListRequest.Mods,
	ApprovedSenders: defaultListRequest.ApprovedSenders,
	Locked:          defaultListRequest.Locked,
}

func TestGetList(t *testing.T) {
	listMock := new(service.IListManagerMock)
	listMock.On("GetList", 1).Return(defaultListResponse)
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/list/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	listMock.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := test.GetExpectedJsonResponse(t, "Successfully fetched list!", defaultListResponse)
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetList401sOnInvalidToken(t *testing.T) {
	authMock := service.NewIAuthManagerMockWithError(util.ErrInvalidToken)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/list/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "mocked error\n", rr.Body.String())
}

func TestGetList400sOnNoParam(t *testing.T) {
	listMock := new(service.IListManagerMock)
	listMock.On("GetList", mock.Anything).Panic("GetList should not have been called")
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/list/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid id provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetList400sOnInvalidParam(t *testing.T) {
	listMock := new(service.IListManagerMock)
	listMock.On("GetList", mock.Anything).Panic("GetList should not have been called")
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/list/?id=a", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid id provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetList404sWhenNoUser(t *testing.T) {
	listMock := service.NewIListManagerMockWithError(util.ErrNoUser)
	listMock.On("GetList", 1).Return(defaultListResponse)
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/list/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetList500sOnGenericError(t *testing.T) {
	listMock := service.NewIListManagerMockWithError(util.Unknown)
	listMock.On("GetList", 1).Return(defaultListResponse)
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/list/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPostList(t *testing.T) {
	expectedList := &model.ListRequest{Name: defaultListRequest.Name, Recipients: defaultListRequest.Recipients,
		Mods: append(defaultListRequest.Mods, 1), ApprovedSenders: defaultListRequest.ApprovedSenders,
		Locked: defaultListRequest.Locked}
	listMock := new(service.IListManagerMock)
	listMock.On("CreateList", expectedList).Return(1)
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	authMock.On("CheckPerms", 1, "CRT_LIST").Return(true)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/list/", defaultListRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	listMock.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusCreated, rr.Code)
	expected := test.GetExpectedJsonResponse(t, "Successfully created list!", IdObject{Id: 1})
	assert.Equal(t, expected, rr.Body.String())
}

func TestPostList401sOnInvalidToken(t *testing.T) {
	authMock := service.NewIAuthManagerMockWithError(util.ErrInvalidToken)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/list/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "mocked error\n", rr.Body.String())
}

func TestPostList400sOnInvalidJson(t *testing.T) {
	listMock := new(service.IListManagerMock)
	listMock.On("CreateList", mock.Anything).Panic("CreateList should not have been called")
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/list/", 1)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid json provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPostList400sOnInvalidObj(t *testing.T) {
	expectedList := &model.ListRequest{Name: defaultListRequest.Name, Recipients: defaultListRequest.Recipients,
		Mods: append(defaultListRequest.Mods, 1), ApprovedSenders: defaultListRequest.ApprovedSenders,
		Locked: defaultListRequest.Locked}
	listMock := service.NewIListManagerMockWithError(util.ErrInvalidObject)
	listMock.On("CreateList", expectedList).Return(1)
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	authMock.On("CheckPerms", 1, "CRT_LIST").Return(true)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/list/", defaultListRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPostList409sOnDuplicateUser(t *testing.T) {
	expectedList := &model.ListRequest{Name: defaultListRequest.Name, Recipients: defaultListRequest.Recipients,
		Mods: append(defaultListRequest.Mods, 1), ApprovedSenders: defaultListRequest.ApprovedSenders,
		Locked: defaultListRequest.Locked}
	listMock := service.NewIListManagerMockWithError(util.ErrListAlreadyExists)
	listMock.On("CreateList", expectedList).Return(1)
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	authMock.On("CheckPerms", 1, "CRT_LIST").Return(true)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/list/", defaultListRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusConflict, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPostList500sOnDbError(t *testing.T) {
	expectedList := &model.ListRequest{Name: defaultListRequest.Name, Recipients: defaultListRequest.Recipients,
		Mods: append(defaultListRequest.Mods, 1), ApprovedSenders: defaultListRequest.ApprovedSenders,
		Locked: defaultListRequest.Locked}
	listMock := service.NewIListManagerMockWithError(util.Unknown)
	listMock.On("CreateList", expectedList).Return(1)
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	authMock.On("CheckPerms", 1, "CRT_LIST").Return(true)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/list/", defaultListRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchList(t *testing.T) {
	listMock := new(service.IListManagerMock)
	listMock.On("UpdateList", 1, defaultListRequest, true).Return()
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	authMock.On("CheckListMods", 1, 1).Return(true)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "PATCH", "/api/list/?id=1", defaultListRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	listMock.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := test.GetExpectedJsonResponse(t, "Successfully patched list!", IdObject{Id: 1})
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchListCallsNoLocked(t *testing.T) {
	expectedList := &model.ListRequest{Name: defaultListRequest.Name, Recipients: defaultListRequest.Recipients,
		Mods: defaultListRequest.Mods, ApprovedSenders: defaultListRequest.ApprovedSenders}
	listMock := new(service.IListManagerMock)
	listMock.On("UpdateList", 1, expectedList, false).Return()
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	authMock.On("CheckListMods", 1, 1).Return(true)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	request := &listRequestWithoutLocked{Name: defaultListRequest.Name, Recipients: defaultListRequest.Recipients,
		Mods: defaultListRequest.Mods, ApprovedSenders: defaultListRequest.ApprovedSenders}
	req := test.CreateHttpRequest(t, "PATCH", "/api/list/?id=1", request)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	listMock.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := test.GetExpectedJsonResponse(t, "Successfully patched list!", IdObject{Id: 1})
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchList401sOnInvalidToken(t *testing.T) {
	authMock := service.NewIAuthManagerMockWithError(util.ErrInvalidToken)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "PATCH", "/api/list/?id=1", defaultListRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "mocked error\n", rr.Body.String())
}

func TestPatchList400sOnNoParam(t *testing.T) {
	listMock := new(service.IListManagerMock)
	listMock.On("UpdateList", mock.Anything, mock.Anything).Panic("UpdateList should not have been called")
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "PATCH", "/api/list/", defaultListRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid id provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchList400sOnInvalidParam(t *testing.T) {
	listMock := new(service.IListManagerMock)
	listMock.On("UpdateList", mock.Anything, mock.Anything).Panic("UpdateList should not have been called")
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "PATCH", "/api/list/?id=a", defaultListRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid id provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchList400sOnInvalidJson(t *testing.T) {
	listMock := new(service.IListManagerMock)
	listMock.On("UpdateList", mock.Anything, mock.Anything).Panic("UpdateList should not have been called")
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "PATCH", "/api/list/?id=1", 1)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid json provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchList400sOnInvalidObj(t *testing.T) {
	listMock := service.NewIListManagerMockWithError(util.ErrInvalidObject)
	listMock.On("UpdateList", 1, defaultListRequest, true).Return()
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	authMock.On("CheckListMods", 1, 1).Return(true)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "PATCH", "/api/list/?id=1", defaultListRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchList409sOnDuplicateName(t *testing.T) {
	listMock := service.NewIListManagerMockWithError(util.ErrListAlreadyExists)
	listMock.On("UpdateList", 1, defaultListRequest, true).Return()
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	authMock.On("CheckListMods", 1, 1).Return(true)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "PATCH", "/api/list/?id=1", defaultListRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusConflict, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchUser404sWhenNoList(t *testing.T) {
	listMock := service.NewIListManagerMockWithError(util.ErrNoList)
	listMock.On("UpdateList", 1, defaultListRequest, true).Return()
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	authMock.On("CheckListMods", 1, 1).Return(true)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "PATCH", "/api/list/?id=1", defaultListRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchList500sOnGenericError(t *testing.T) {
	listMock := service.NewIListManagerMockWithError(util.Unknown)
	listMock.On("UpdateList", 1, defaultListRequest, true).Return()
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	authMock.On("CheckListMods", 1, 1).Return(true)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "PATCH", "/api/list/?id=1", defaultListRequest)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestDeleteList(t *testing.T) {
	listMock := new(service.IListManagerMock)
	listMock.On("DeleteList", 1).Return()
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	authMock.On("CheckListMods", 1, 1).Return(true)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "DELETE", "/api/list/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	listMock.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := test.GetExpectedJsonResponse(t, "Successfully deleted list!", IdObject{Id: 1})
	assert.Equal(t, expected, rr.Body.String())
}

func TestDeleteList401sOnInvalidToken(t *testing.T) {
	authMock := service.NewIAuthManagerMockWithError(util.ErrInvalidToken)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "DELETE", "/api/list/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "mocked error\n", rr.Body.String())
}

func TestDeleteList400sOnNoParam(t *testing.T) {
	listMock := new(service.IListManagerMock)
	listMock.On("DeleteList", mock.Anything).Panic("DeleteUser should not have been called")
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "DELETE", "/api/list/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid id provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestDeleteList404sOnInvalidParam(t *testing.T) {
	listMock := new(service.IListManagerMock)
	listMock.On("DeleteList", mock.Anything).Panic("DeleteUser should not have been called")
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "DELETE", "/api/list/?id=a", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid id provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestDeleteList404sWhenNoUser(t *testing.T) {
	listMock := service.NewIListManagerMockWithError(util.ErrNoList)
	listMock.On("DeleteList", 1).Return()
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	authMock.On("CheckListMods", 1, 1).Return(true)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "DELETE", "/api/list/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestDeleteList500sOnGenericError(t *testing.T) {
	listMock := service.NewIListManagerMockWithError(util.Unknown)
	listMock.On("DeleteList", 1).Return()
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	authMock.On("CheckListMods", 1, 1).Return(true)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "DELETE", "/api/list/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetLists(t *testing.T) {
	listMock := new(service.IListManagerMock)
	listMock.On("GetAllLists").Return(&[]*model.ListResponse{defaultListResponse})
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/lists/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	listMock.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := test.GetExpectedJsonResponse(t, "Successfully fetched lists!", []*model.ListResponse{defaultListResponse})
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetLists401sOnInvalidToken(t *testing.T) {
	authMock := service.NewIAuthManagerMockWithError(util.ErrInvalidToken)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/lists/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "mocked error\n", rr.Body.String())
}

func TestGetLists500sOnGenericError(t *testing.T) {
	listMock := service.NewIListManagerMockWithError(util.Unknown)
	listMock.On("GetAllLists").Return(&[]*model.ListResponse{defaultListResponse})
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addListHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/lists/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "mocked error\n"
	assert.Equal(t, expected, rr.Body.String())
}
