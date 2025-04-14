package mock

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
)

type IDbMock struct {
	mock.Mock
	db.IDb
	error *db.Error
}

func NewIDbMockWithError(code db.ErrorCode) *IDbMock {
	return &IDbMock{error: &db.Error{Code: code}}
}

func (mock *IDbMock) GetList(id int) (*model.List, *db.Error) {
	args := mock.Called(id)
	return args.Get(0).(*model.List), mock.error
}

func (mock *IDbMock) CreateList(name string, recipients []string) (int, *db.Error) {
	args := mock.Called(name, recipients)
	return args.Int(0), mock.error
}

func (mock *IDbMock) PatchList(id int, name string, recipients []string) *db.Error {
	mock.Called(id, name, recipients)
	return mock.error
}

func (mock *IDbMock) DeleteList(id int) *db.Error {
	mock.Called(id)
	return mock.error
}

func (mock *IDbMock) GetAllLists() (*[]*model.List, *db.Error) {
	args := mock.Called()
	return args.Get(0).(*[]*model.List), mock.error
}

func (mock *IDbMock) GetRecipientsFromListName(name string) ([]string, error) {
	args := mock.Called(name)
	if mock.error != nil {
		return nil, errors.New("mocked error")
	}
	return args.Get(0).([]string), nil
}
