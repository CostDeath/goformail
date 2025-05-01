package service

import (
	"github.com/stretchr/testify/mock"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
)

type IAuthManagerMock struct {
	mock.Mock
	IAuthManager
	error *util.Error
}

func NewIAuthManagerMockWithError(code util.ErrorCode) *IAuthManagerMock {
	return &IAuthManagerMock{error: &util.Error{Code: code, Message: "mocked error"}}
}

func (mock *IAuthManagerMock) Login(creds *model.LoginRequest) (string, int, *util.Error) {
	args := mock.Called(creds)
	return args.String(0), args.Int(1), mock.error
}

func (mock *IAuthManagerMock) CheckTokenValidity(tokenStr string) (int, *util.Error) {
	args := mock.Called(tokenStr)
	return args.Int(0), mock.error
}

func (mock *IAuthManagerMock) CheckPerms(id int, required string) (bool, *util.Error) {
	args := mock.Called(id, required)
	return args.Bool(0), mock.error
}

func (mock *IAuthManagerMock) CheckUserPerms(id int, action string, required []string) (bool, *util.Error) {
	args := mock.Called(id, action, required)
	return args.Bool(0), mock.error
}

func (mock *IAuthManagerMock) CheckListMods(id int, listId int) (bool, *util.Error) {
	args := mock.Called(id, listId)
	return args.Bool(0), mock.error
}
