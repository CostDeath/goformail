package cli

import (
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"testing"
)

func TestCreateList(t *testing.T) {
	getMock := func(db db.IDb) service.IListManager {
		mockList := new(service.IListManagerMock)
		mockList.On("CreateList", 1, &model.ListRequest{
			Name: "", Recipients: []string{}, ApprovedSenders: []string{}, Mods: []int64{}, Locked: false},
		).Return(nil)
		return mockList
	}

	createList([]string{"1"}, &db.Db{}, getMock)
}

func TestCreateListAllSet(t *testing.T) {
	getMock := func(db db.IDb) service.IListManager {
		mockList := new(service.IListManagerMock)
		mockList.On("CreateList", 1, &model.ListRequest{
			Name: "list", Recipients: []string{"rcpt@domain.tld", "rcpt2@domain.tld"},
			ApprovedSenders: []string{"s@domain.tld"}, Mods: []int64{1, 2}, Locked: false},
		).Return(nil)
		return mockList
	}

	createList([]string{"1", "--name=list", "--recipients=rcpt@domain.tld,rcpt2@domain.tld", "approved=s@domain.tld",
		"mod=1", "mod=2", "locked=false"}, &db.Db{}, getMock)
}

func TestCreateUser(t *testing.T) {
	getMock := func(db db.IDb) service.IUserManager {
		mockUser := new(service.IUserManagerMock)
		mockUser.On("CreateUser", 1, &model.UserRequest{
			Email: "", Password: "", Permissions: []string{},
		}).Return(nil)
		return mockUser
	}

	createUser([]string{"1"}, &db.Db{}, getMock)
}

func TestCreateUserAllSet(t *testing.T) {
	getMock := func(db db.IDb) service.IUserManager {
		mockUser := new(service.IUserManagerMock)
		mockUser.On("CreateUser", 1, &model.UserRequest{
			Email: "user@domain.tld", Password: "password", Permissions: []string{"ADMIN", "CRT_LIST", "MOD_LIST",
				"CRT_USER", "MOD_USER"},
		}).Return(nil)
		return mockUser
	}

	createUser([]string{"1", "--email=user@domain.tld", "--password=password",
		"--permission=ADMIN,CRT_LIST,MOD_LIST,CRT_USER,MOD_USER"}, &db.Db{}, getMock)
}
