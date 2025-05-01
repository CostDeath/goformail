package db

import (
	"github.com/stretchr/testify/mock"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
)

type IDbMock struct {
	mock.Mock
	IDb
	error *Error
}

func NewIDbMockWithError(code ErrorCode) *IDbMock {
	return &IDbMock{error: &Error{Code: code}}
}

func (mock *IDbMock) GetList(id int) (*model.ListResponse, *Error) {
	args := mock.Called(id)
	return args.Get(0).(*model.ListResponse), mock.error
}

func (mock *IDbMock) CreateList(list *model.ListRequest) (int, *Error) {
	args := mock.Called(list)
	return args.Int(0), mock.error
}

func (mock *IDbMock) PatchList(id int, list *model.ListRequest, override *model.ListOverrides) *Error {
	mock.Called(id, list, override)
	return mock.error
}

func (mock *IDbMock) DeleteList(id int) *Error {
	mock.Called(id)
	return mock.error
}

func (mock *IDbMock) GetAllLists() (*[]*model.ListResponse, *Error) {
	args := mock.Called()
	return args.Get(0).(*[]*model.ListResponse), mock.error
}

func (mock *IDbMock) GetApprovalFromListName(sender string, name string) (int, bool, *Error) {
	args := mock.Called(sender, name)
	return args.Int(0), args.Bool(1), mock.error
}

func (mock *IDbMock) GetUser(id int) (*model.UserResponse, *Error) {
	args := mock.Called(id)
	return args.Get(0).(*model.UserResponse), mock.error
}

func (mock *IDbMock) CreateUser(user *model.UserRequest, hash string) (int, *Error) {
	args := mock.Called(user, hash)
	return args.Int(0), mock.error
}

func (mock *IDbMock) UpdateUser(id int, user *model.UserRequest, overridePerms bool) *Error {
	mock.Called(id, user, overridePerms)
	return mock.error
}

func (mock *IDbMock) DeleteUser(id int) *Error {
	mock.Called(id)
	return mock.error
}

func (mock *IDbMock) GetAllUsers() (*[]*model.UserResponse, *Error) {
	args := mock.Called()
	return args.Get(0).(*[]*model.UserResponse), mock.error
}

func (mock *IDbMock) GetUserPassword(email string) (int, string, *Error) {
	args := mock.Called(email)
	return args.Int(0), args.String(1), mock.error
}

func (mock *IDbMock) UserExists(id int) (bool, *Error) {
	args := mock.Called(id)
	return args.Bool(0), mock.error
}

func (mock *IDbMock) UsersExist(ids []int64) ([]int64, *Error) {
	args := mock.Called(ids)
	return args.Get(0).([]int64), mock.error
}

func (mock *IDbMock) GetUserPerms(id int) ([]string, *Error) {
	args := mock.Called(id)
	return args.Get(0).([]string), mock.error
}

func (mock *IDbMock) GetUserPermsAndModStatus(id int, listId int) ([]string, bool, *Error) {
	args := mock.Called(id, listId)
	return args.Get(0).([]string), args.Bool(1), mock.error
}

func (mock *IDbMock) GetAllReadyEmails() (*[]model.Email, *Error) {
	args := mock.Called()
	return args.Get(0).(*[]model.Email), mock.error
}

func (mock *IDbMock) AddEmail(email *model.Email) *Error {
	mock.Called(email)
	return mock.error
}

func (mock *IDbMock) SetEmailAsSent(id int) *Error {
	mock.Called(id)
	return mock.error
}

func (mock *IDbMock) SetEmailRetry(email *model.Email) *Error {
	mock.Called(email)
	return mock.error
}
