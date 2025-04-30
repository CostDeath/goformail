package service

import (
	"github.com/stretchr/testify/mock"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
)

type IListManagerMock struct {
	mock.Mock
	IListManager
	error *util.Error
}

func NewIListManagerMockWithError(code util.ErrorCode) *IListManagerMock {
	return &IListManagerMock{error: &util.Error{Code: code, Message: "mocked error"}}
}

func (mock *IListManagerMock) GetList(id int) (*model.ListResponse, *util.Error) {
	args := mock.Called(id)
	return args.Get(0).(*model.ListResponse), mock.error
}

func (mock *IListManagerMock) CreateList(list *model.ListRequest) (int, *util.Error) {
	args := mock.Called(list)
	return args.Int(0), mock.error
}

func (mock *IListManagerMock) UpdateList(id int, list *model.ListRequest) *util.Error {
	mock.Called(id, list)
	return mock.error
}

func (mock *IListManagerMock) DeleteList(id int) *util.Error {
	mock.Called(id)
	return mock.error
}

func (mock *IListManagerMock) GetAllLists() (*[]*model.ListResponse, *util.Error) {
	args := mock.Called()
	return args.Get(0).(*[]*model.ListResponse), mock.error
}
