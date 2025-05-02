package rest

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

var list = 1
var defaultReqs = &model.EmailReqs{Offset: 10, List: &list, Archived: true, Exhausted: true, PendingApproval: true}
var defaultEmail = &model.Email{Id: 1, Content: "content"}

func TestGetEmail(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetEmail", 1).Return(defaultEmail)
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{db: dbMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addEmailHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/email/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	dbMock.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := test.GetExpectedJsonResponse(t, "Successfully fetched email!", defaultEmail)
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetEmail401sOnInvalidToken(t *testing.T) {
	authMock := service.NewIAuthManagerMockWithError(util.ErrInvalidToken)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{auth: authMock, mux: new(http.ServeMux)}
	ctrl.addEmailHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/email/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "mocked error\n", rr.Body.String())
}

func TestGetEmail400sOnNoParam(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetEmail", mock.Anything).Panic("GetEmail should not have been called")
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{db: dbMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addEmailHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/email/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid id provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetEmail400sOnInvalidParam(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetEmail", mock.Anything).Panic("GetEmail should not have been called")
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{db: dbMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addEmailHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/email/?id=a", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid id provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetList404sWhenNoEmail(t *testing.T) {
	dbMock := db.NewIDbMockWithError(db.ErrNoRows)
	dbMock.On("GetEmail", 1).Return(defaultEmail)
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{db: dbMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addEmailHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/email/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "Could not find a email with id '1'\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetEmail500sOnGenericError(t *testing.T) {
	dbMock := db.NewIDbMockWithError(db.Unknown)
	dbMock.On("GetEmail", 1).Return(defaultEmail)
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{db: dbMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addEmailHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "GET", "/api/email/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "Unknown error occurred\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetEmails(t *testing.T) {
	resp := &model.EmailResponse{Offset: 20, Emails: []model.Email{{Id: 1, Content: "content"}}}
	dbMock := new(db.IDbMock)
	dbMock.On("GetAllEmails", defaultReqs).Return(resp)
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{db: dbMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addEmailHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/emails/", defaultReqs)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	dbMock.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := test.GetExpectedJsonResponse(t, "Successfully fetched emails!", resp)
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetEmails401sOnInvalidToken(t *testing.T) {
	authMock := service.NewIAuthManagerMockWithError(util.ErrInvalidToken)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{auth: authMock, mux: new(http.ServeMux)}
	ctrl.addEmailHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/emails/", defaultReqs)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "mocked error\n", rr.Body.String())
}

func TestGetEmails400sOnInvalidJson(t *testing.T) {
	listMock := new(service.IListManagerMock)
	listMock.On("CreateList", mock.Anything).Panic("CreateList should not have been called")
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{list: listMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addEmailHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/emails/", 1)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid json provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetEmails500sOnDbError(t *testing.T) {
	resp := &model.EmailResponse{Offset: 20, Emails: []model.Email{{Id: 1, Content: "content"}}}
	dbMock := db.NewIDbMockWithError(db.Unknown)
	dbMock.On("GetAllEmails", defaultReqs).Return(resp)
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{db: dbMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addEmailHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/emails/", defaultReqs)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "Unknown error occurred\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestApproveEmail(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetEmailList", 1).Return(2)
	dbMock.On("SetEmailAsApproved", 1).Return(nil)
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(3)
	authMock.On("CheckListMods", 3, 2).Return(true)
	ctrl := &Controller{db: dbMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addEmailHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/emails/approve/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check db was called with correct args
	dbMock.AssertExpectations(t)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := test.GetExpectedJsonResponse(t, "Successfully approved email!", IdObject{Id: 1})
	assert.Equal(t, expected, rr.Body.String())
}

func TestApproveEmail401sOnInvalidToken(t *testing.T) {
	authMock := service.NewIAuthManagerMockWithError(util.ErrInvalidToken)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{auth: authMock, mux: new(http.ServeMux)}
	ctrl.addEmailHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/emails/approve/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "mocked error\n", rr.Body.String())
}

func TestApproveEmail400sOnNoParam(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetEmailList", mock.Anything).Panic("GetEmailList should not have been called")
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{db: dbMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addEmailHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/emails/approve/", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid id provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestApproveEmail400sOnInvalidParam(t *testing.T) {
	dbMock := new(db.IDbMock)
	dbMock.On("GetEmailList", mock.Anything).Panic("GetEmailList should not have been called")
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	ctrl := &Controller{db: dbMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addEmailHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/emails/approve/?id=a", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expected := "Invalid object: Invalid id provided\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestApproveEmail404sOnNoEmailError(t *testing.T) {
	dbMock := db.NewIDbMockWithError(db.ErrNoRows)
	dbMock.On("GetEmailList", 1).Return(0)
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	authMock.On("CheckListMods", 1, 1).Return(true)
	ctrl := &Controller{db: dbMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addEmailHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/emails/approve/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)
	expected := "Could not find a email with id '1'\n"
	assert.Equal(t, expected, rr.Body.String())
}

func TestApproveEmail500sOnGenericError(t *testing.T) {
	dbMock := db.NewIDbMockWithError(db.Unknown)
	dbMock.On("GetEmailList", 1).Return(0)
	authMock := new(service.IAuthManagerMock)
	authMock.On("CheckTokenValidity", "Bearer token").Return(1)
	authMock.On("CheckListMods", 1, 1).Return(true)
	ctrl := &Controller{db: dbMock, auth: authMock, mux: new(http.ServeMux)}
	ctrl.addEmailHandlers()

	// Mock the request
	req := test.CreateHttpRequest(t, "POST", "/api/emails/approve/?id=1", nil)
	rr := httptest.NewRecorder()
	ctrl.mux.ServeHTTP(rr, req)

	// Check the response is what we expect.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expected := "Unknown error occurred\n"
	assert.Equal(t, expected, rr.Body.String())
}
