package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"testing"
)

var defaultUserRequest = &model.UserRequest{
	Email: "example@domain.tld", Password: "pass", Permissions: []string{"ADMIN", "CRT_LIST"},
}
var defaultUserRequestUpper = &model.UserRequest{
	Email: "EXAMPLE@domain.tld", Password: defaultUserRequest.Password, Permissions: defaultUserRequest.Permissions,
}
var defaultUserResponse = &model.UserResponse{
	Id: 1, Email: defaultUserRequest.Email, Permissions: defaultUserRequest.Permissions,
}

func TestGetUser(t *testing.T) {
	mockObj := new(db.IDbMock)
	mockObj.On("GetUser", 1).Return(defaultUserResponse)
	man := UserManager{db: mockObj}
	actual, err := man.GetUser(1)

	mockObj.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, defaultUserResponse, actual)
}

func TestGetUserReturnsNoUserError(t *testing.T) {
	mockObj := db.NewIDbMockWithError(db.ErrNoRows)
	mockObj.On("GetUser", 1).Return(defaultUserResponse)
	man := UserManager{db: mockObj}
	actual, err := man.GetUser(1)

	mockObj.AssertExpectations(t)
	assert.Equal(t, util.NewNoUserError(1, nil), err)
	assert.Nil(t, actual)
}

func TestGetUserReturnsGenericError(t *testing.T) {
	mockObj := db.NewIDbMockWithError(db.Unknown)
	mockObj.On("GetUser", 1).Return(defaultUserResponse)
	man := UserManager{db: mockObj}
	actual, err := man.GetUser(1)

	mockObj.AssertExpectations(t)
	assert.Equal(t, util.NewGenericError(nil), err)
	assert.Nil(t, actual)
}

func TestCreateUser(t *testing.T) {
	mockObj := new(db.IDbMock)
	mockObj.On("CreateUser", defaultUserRequest, "hash", "salt").Return(1)
	man := UserManager{db: mockObj}
	id, err := man.CreateUser(defaultUserRequestUpper)

	mockObj.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, 1, id)
}

func TestCreateUserReturnsInvalidObjectErrorOnMissingFields(t *testing.T) {
	mockObj := new(db.IDbMock)
	mockObj.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
		Panic("CreateUser should not have been called")
	man := UserManager{db: mockObj}
	id, err := man.CreateUser(&model.UserRequest{})

	expected := util.NewInvalidObjectError("Missing field(s) in user: Email, Password, Permissions", nil)
	assert.Equal(t, expected, err)
	assert.Equal(t, 0, id)
}

func TestCreateUserReturnsInvalidObjectErrorOnInvalidEmail(t *testing.T) {
	mockObj := new(db.IDbMock)
	mockObj.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
		Panic("CreateUser should not have been called")
	man := UserManager{db: mockObj}
	id, err := man.CreateUser(&model.UserRequest{Email: "invalid", Password: "pass", Permissions: []string{}})

	expected := util.NewInvalidObjectError("Invalid email address 'invalid'", nil)
	assert.Equal(t, expected, err)
	assert.Equal(t, 0, id)
}

func TestCreateUserReturnsInvalidObjectErrorOnInvalidPerms(t *testing.T) {
	mockObj := new(db.IDbMock)
	mockObj.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
		Panic("CreateUser should not have been called")
	man := UserManager{db: mockObj}
	id, err := man.CreateUser(&model.UserRequest{
		Email: "example@domain.tld", Password: "pass", Permissions: []string{"invalid"}})

	msg := "Missing or duplicate permission. Valid permissions- ADMIN, CRT_LIST, MOD_LIST, CRT_USER, MOD_USER"
	expected := util.NewInvalidObjectError(msg, nil)
	assert.Equal(t, expected, err)
	assert.Equal(t, 0, id)
}

func TestCreateUserReturnsUserAlreadyExistsError(t *testing.T) {
	mockObj := db.NewIDbMockWithError(db.ErrDuplicate)
	mockObj.On("CreateUser", defaultUserRequest, "hash", "salt").Return(0)
	man := UserManager{db: mockObj}
	id, err := man.CreateUser(defaultUserRequest)

	assert.Equal(t, util.NewUserAlreadyExistsError(defaultUserRequest.Email, nil), err)
	assert.Equal(t, 0, id)
}

func TestCreateUserReturnsGenericError(t *testing.T) {
	mockObj := db.NewIDbMockWithError(db.Unknown)
	mockObj.On("CreateUser", defaultUserRequest, "hash", "salt").Return(0)
	man := UserManager{db: mockObj}
	id, err := man.CreateUser(defaultUserRequest)

	assert.Equal(t, util.NewGenericError(nil), err)
	assert.Equal(t, 0, id)
}

func TestUpdateUser(t *testing.T) {
	mockObj := new(db.IDbMock)
	mockObj.On("UpdateUser", 1, defaultUserRequest).Return()
	man := UserManager{db: mockObj}
	err := man.UpdateUser(1, defaultUserRequestUpper)

	mockObj.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestUpdateUserReturnsInvalidObjectErrorOnInvalidEmail(t *testing.T) {
	mockObj := new(db.IDbMock)
	mockObj.On("UpdateUser", mock.Anything, defaultUserRequest).
		Panic("UpdateUser should not have been called")
	man := UserManager{db: mockObj}
	err := man.UpdateUser(1, &model.UserRequest{Email: "invalid"})

	expected := util.NewInvalidObjectError("Invalid email address 'invalid'", nil)
	assert.Equal(t, expected, err)
}

func TestUpdateUserReturnsInvalidObjectErrorOnInvalidPerms(t *testing.T) {
	mockObj := new(db.IDbMock)
	mockObj.On("UpdateUser", mock.Anything, defaultUserRequest).
		Panic("UpdateUser should not have been called")
	man := UserManager{db: mockObj}
	err := man.UpdateUser(1, &model.UserRequest{Permissions: []string{"invalid"}})

	msg := "Missing or duplicate permission. Valid permissions- ADMIN, CRT_LIST, MOD_LIST, CRT_USER, MOD_USER"
	expected := util.NewInvalidObjectError(msg, nil)
	assert.Equal(t, expected, err)
}

func TestUpdateUserReturnsUserAlreadyExistsError(t *testing.T) {
	mockObj := db.NewIDbMockWithError(db.ErrDuplicate)
	mockObj.On("UpdateUser", 1, defaultUserRequest).Return()
	man := UserManager{db: mockObj}
	err := man.UpdateUser(1, defaultUserRequest)

	assert.Equal(t, util.NewUserAlreadyExistsError(defaultUserRequest.Email, nil), err)
}

func TestUpdateUserReturnsNoUserError(t *testing.T) {
	mockObj := db.NewIDbMockWithError(db.ErrNoRows)
	mockObj.On("UpdateUser", 1, defaultUserRequest).Return()
	man := UserManager{db: mockObj}
	err := man.UpdateUser(1, defaultUserRequest)

	mockObj.AssertExpectations(t)
	assert.Equal(t, util.NewNoUserError(1, nil), err)
}

func TestUpdateUserReturnsGenericError(t *testing.T) {
	mockObj := db.NewIDbMockWithError(db.Unknown)
	mockObj.On("UpdateUser", 1, defaultUserRequest).Return()
	man := UserManager{db: mockObj}
	err := man.UpdateUser(1, defaultUserRequest)

	assert.Equal(t, util.NewGenericError(nil), err)
}

func TestDeleteUser(t *testing.T) {
	mockObj := new(db.IDbMock)
	mockObj.On("DeleteUser", 1).Return()
	man := UserManager{db: mockObj}
	err := man.DeleteUser(1)

	mockObj.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestDeleteUserReturnsNoUserError(t *testing.T) {
	mockObj := db.NewIDbMockWithError(db.ErrNoRows)
	mockObj.On("DeleteUser", 1).Return()
	man := UserManager{db: mockObj}
	err := man.DeleteUser(1)

	mockObj.AssertExpectations(t)
	assert.Equal(t, util.NewNoUserError(1, nil), err)
}

func TestDeleteUserReturnsGenericError(t *testing.T) {
	mockObj := db.NewIDbMockWithError(db.Unknown)
	mockObj.On("DeleteUser", 1).Return()
	man := UserManager{db: mockObj}
	err := man.DeleteUser(1)

	assert.Equal(t, util.NewGenericError(nil), err)
}

func TestGetAllUsers(t *testing.T) {
	mockObj := new(db.IDbMock)
	expected := &[]*model.UserResponse{defaultUserResponse}
	mockObj.On("GetAllUsers").Return(expected)
	man := UserManager{db: mockObj}
	actual, err := man.GetAllUsers()

	mockObj.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestGetLists500sOnDbError(t *testing.T) {
	mockObj := db.NewIDbMockWithError(db.Unknown)
	mockObj.On("GetAllUsers").Return(&[]*model.UserResponse{defaultUserResponse})
	man := UserManager{db: mockObj}
	actual, err := man.GetAllUsers()

	assert.Equal(t, util.NewGenericError(nil), err)
	assert.Nil(t, actual)
}
