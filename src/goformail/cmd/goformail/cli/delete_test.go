package cli

import (
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"testing"
)

func TestDeleteList(t *testing.T) {
	getMock := func(db db.IDb) service.IListManager {
		mockList := new(service.IListManagerMock)
		mockList.On("DeleteList", 1).Return(nil)
		return mockList
	}

	createList([]string{"1"}, &db.Db{}, getMock)
}

func TestDeleteUser(t *testing.T) {
	getMock := func(db db.IDb) service.IUserManager {
		mockUser := new(service.IUserManagerMock)
		mockUser.On("DeleteUser", 1).Return(nil)
		return mockUser
	}

	createUser([]string{"1"}, &db.Db{}, getMock)
}
