package service

import (
	"github.com/stretchr/testify/mock"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
)

type IUserManagerMock struct {
	mock.Mock
	IUserManager
	error *util.Error
}

func NewIUserManagerMockWithError(code util.ErrorCode) *IUserManagerMock {
	return &IUserManagerMock{error: &util.Error{Code: code, Message: "mocked error"}}
}

func (mock *IUserManagerMock) GetUser(id int) (*model.UserResponse, *util.Error) {
	args := mock.Called(id)
	return args.Get(0).(*model.UserResponse), mock.error
}

func (mock *IUserManagerMock) CreateUser(user *model.UserRequest) (int, *util.Error) {
	args := mock.Called(user)
	return args.Int(0), mock.error
}

func (mock *IUserManagerMock) UpdateUser(id int, user *model.UserRequest) *util.Error {
	mock.Called(id, user)
	return mock.error
}

func (mock *IUserManagerMock) DeleteUser(id int) *util.Error {
	mock.Called(id)
	return mock.error
}

func (mock *IUserManagerMock) GetAllUsers() (*[]*model.UserResponse, *util.Error) {
	args := mock.Called()
	return args.Get(0).(*[]*model.UserResponse), mock.error
}
