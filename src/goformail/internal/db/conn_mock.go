package db

import (
	"errors"
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

func (mock *IDbMock) GetList(id int) (*model.List, *Error) {
	args := mock.Called(id)
	return args.Get(0).(*model.List), mock.error
}

func (mock *IDbMock) CreateList(list *model.List) (int, *Error) {
	args := mock.Called(list)
	return args.Int(0), mock.error
}

func (mock *IDbMock) PatchList(id int, list *model.List) *Error {
	mock.Called(id, list)
	return mock.error
}

func (mock *IDbMock) DeleteList(id int) *Error {
	mock.Called(id)
	return mock.error
}

func (mock *IDbMock) GetAllLists() (*[]*model.ListWithId, *Error) {
	args := mock.Called()
	return args.Get(0).(*[]*model.ListWithId), mock.error
}

func (mock *IDbMock) GetRecipientsFromListName(name string) ([]string, error) {
	args := mock.Called(name)
	if mock.error != nil {
		return nil, errors.New("mocked error")
	}
	return args.Get(0).([]string), nil
}
