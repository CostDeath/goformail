package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/rest/util"
	util2 "gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"net/http"
	"net/http/httptest"
	"testing"
)

var defaultList = &model.List{Name: "name", Recipients: []string{"example@domain.tld"}}
var defaultListWithId = &model.ListWithId{Id: 1, List: defaultList}

func listsCleanUp() {
	http.DefaultServeMux = new(http.ServeMux)
}

func TestGetList(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := new(db.IDbMock)
	mockObj.On("GetList", 1).Return(defaultList)
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "GET", "/api/list/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	mockObj.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := getExpectedListResponse(t, "Successfully fetched list!", defaultList)
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetList404sOnNoParam(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := new(db.IDbMock)
	mockObj.On("GetList", mock.Anything).Panic("GetList should not have been called")
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "GET", "/api/list/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "List not found\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetList404sOnInvalidParam(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := new(db.IDbMock)
	mockObj.On("GetList", mock.Anything).Panic("GetList should not have been called")
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "GET", "/api/list/?id=a", defaultList)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "List not found\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetList404sWhenNoSuchId(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := db.NewIDbMockWithError(db.ErrNoRows)
	mockObj.On("GetList", 1).Return(defaultList)
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "GET", "/api/list/?id=1", defaultList)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "List not found\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetList500sOnDbError(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := db.NewIDbMockWithError(db.Unknown)
	mockObj.On("GetList", 1).Return(defaultList)
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "GET", "/api/list/?id=1", defaultList)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "An error occurred relating to the database\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPostList(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := new(db.IDbMock)
	mockObj.On("CreateList", defaultList).Return(1)
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "POST", "/api/list/", defaultList)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	mockObj.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusCreated, rr.Code)
	expected := getExpectedListResponse(t, "Successfully created list!", util.IdObject{Id: 1})
	assert.Equal(t, expected, rr.Body.String())
}

func TestPostList400sOnMissingField(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := new(db.IDbMock)
	mockObj.On("CreateList", mock.Anything).
		Panic("CreateList should not have been called")
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	list := &model.List{Name: "name"}
	req := createListRequest(t, "POST", "/api/list/", list)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid payload\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPostList400sOnInvalidPayload(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := new(db.IDbMock)
	mockObj.On("CreateList", mock.Anything).
		Panic("CreateList should not have been called")
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "POST", "/api/list/", 1)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid payload\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPostList409sOnDuplicateName(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := db.NewIDbMockWithError(db.ErrDuplicate)
	mockObj.On("CreateList", defaultList).Return(0)
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "POST", "/api/list/", defaultList)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusConflict, rr.Code)
	expected := "A list with this name already exists\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPostList500sOnDbError(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := db.NewIDbMockWithError(db.Unknown)
	mockObj.On("CreateList", defaultList).Return(0)
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "POST", "/api/list/", defaultList)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "An error occurred relating to the database\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchList(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := new(db.IDbMock)
	mockObj.On("PatchList", 1, defaultList).Return(1)
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "PATCH", "/api/list/?id=1", defaultList)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	mockObj.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := getExpectedListResponse(t, "Successfully patched list!", util.IdObject{Id: 1})
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchList404sOnNoParam(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := new(db.IDbMock)
	mockObj.On("PatchList", mock.Anything, mock.Anything).
		Panic("PatchList should not have been called")
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "PATCH", "/api/list/", defaultList)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "List not found\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchList404sOnInvalidParam(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := new(db.IDbMock)
	mockObj.On("PatchList", mock.Anything, mock.Anything).
		Panic("PatchList should not have been called")
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "PATCH", "/api/list/?id=a", defaultList)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "List not found\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchList400sOnInvalidPayload(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := new(db.IDbMock)
	mockObj.On("PatchList", mock.Anything, mock.Anything).
		Panic("PatchList should not have been called")
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "PATCH", "/api/list/?id=1", 1)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid payload\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchList409sOnDuplicateName(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := db.NewIDbMockWithError(db.ErrDuplicate)
	mockObj.On("PatchList", 1, defaultList).Return(0)
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "PATCH", "/api/list/?id=1", defaultList)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusConflict, rr.Code)
	expected := "A list with this name already exists\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchList404sWhenNoSuchId(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := db.NewIDbMockWithError(db.ErrNoRows)
	mockObj.On("PatchList", 1, defaultList).Return(0)
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "PATCH", "/api/list/?id=1", defaultList)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "List not found\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestPatchList500sOnDbError(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := db.NewIDbMockWithError(db.Unknown)
	mockObj.On("PatchList", 1, defaultList).Return(0)
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "PATCH", "/api/list/?id=1", defaultList)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "An error occurred relating to the database\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestDeleteList(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := new(db.IDbMock)
	mockObj.On("DeleteList", 1).Return()
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "DELETE", "/api/list/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	mockObj.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := getExpectedListResponse(t, "Successfully deleted list!", util.IdObject{Id: 1})
	assert.Equal(t, expected, rr.Body.String())
}

func TestDeleteList404sOnNoParam(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := new(db.IDbMock)
	mockObj.On("DeleteList", mock.Anything).Panic("DeleteList should not have been called")
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "DELETE", "/api/list/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "List not found\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestDeleteList404sOnInvalidParam(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := new(db.IDbMock)
	mockObj.On("DeleteList", mock.Anything).Panic("DeleteList should not have been called")
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "DELETE", "/api/list/?id=a", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "List not found\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestDeleteList404sWhenNoSuchId(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := db.NewIDbMockWithError(db.ErrNoRows)
	mockObj.On("DeleteList", 1).Return()
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "DELETE", "/api/list/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "List not found\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestDeleteList500sOnDbError(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := db.NewIDbMockWithError(db.Unknown)
	mockObj.On("DeleteList", 1).Return()
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "DELETE", "/api/list/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "An error occurred relating to the database\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetLists(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := new(db.IDbMock)
	mockObj.On("GetAllLists").Return(&[]*model.ListWithId{defaultListWithId})
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "GET", "/api/lists/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	mockObj.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := getExpectedListResponse(t, "Successfully fetched lists!", []*model.ListWithId{defaultListWithId})
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetListsList404sWhenNoLists(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := db.NewIDbMockWithError(db.ErrNoRows)
	mockObj.On("GetAllLists").Return(&[]*model.ListWithId{defaultListWithId})
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "GET", "/api/lists/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "List not found\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetLists500sOnDbError(t *testing.T) {
	t.Cleanup(listsCleanUp)
	mockObj := db.NewIDbMockWithError(db.Unknown)
	mockObj.On("GetAllLists").Return(&[]*model.ListWithId{defaultListWithId})
	ctrl := NewController(util2.MockConfigs, mockObj)
	ctrl.addListHandlers()

	// Mock the request
	req := createListRequest(t, "GET", "/api/lists/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "An error occurred relating to the database\n"
	assert.Equal(t, expected, rr.Body.String())
}

func createListRequest(t *testing.T, method string, uri string, body interface{}) *http.Request {
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	return req
}

func getExpectedListResponse(t *testing.T, msg string, data interface{}) string {
	jsonData, err := json.Marshal(data)
	require.NoError(t, err)
	return fmt.Sprintf("{\"message\":\"%s\",\"data\":%s}", msg, jsonData)
}
