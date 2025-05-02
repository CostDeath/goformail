package cli

import (
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"testing"
)

func TestUpdateList(t *testing.T) {
	var nilStrList []string
	var nilIntList []int64
	getMock := func(db db.IDb) service.IListManager {
		mockList := new(service.IListManagerMock)
		mockList.On("UpdateList", 1, &model.ListRequest{
			Name: "", Recipients: nilStrList, ApprovedSenders: nilStrList, Mods: nilIntList, Locked: false}, false,
		).Return(nil)
		return mockList
	}

	updateList([]string{"1"}, &db.Db{}, getMock)
}

func TestUpdateListAllSet(t *testing.T) {
	getMock := func(db db.IDb) service.IListManager {
		mockList := new(service.IListManagerMock)
		mockList.On("UpdateList", 1, &model.ListRequest{
			Name: "list", Recipients: []string{"rcpt@domain.tld", "rcpt2@domain.tld"},
			ApprovedSenders: []string{"s@domain.tld"}, Mods: []int64{1, 2}, Locked: false}, true,
		).Return(nil)
		return mockList
	}

	updateList([]string{"1", "--name=list", "--recipient=rcpt@domain.tld,rcpt2@domain.tld", "approved=s@domain.tld",
		"mod=1", "mod=2", "locked=false"}, &db.Db{}, getMock)
}

func TestUpdateUser(t *testing.T) {
	getMock := func(db db.IDb) service.IUserManager {
		var nilPerms []string
		mockUser := new(service.IUserManagerMock)
		mockUser.On("UpdateUser", 1, &model.UserRequest{
			Email: "", Permissions: nilPerms,
		}).Return(nil)
		return mockUser
	}

	updateUser([]string{"1"}, &db.Db{}, getMock)
}

func TestUpdateUserAllSet(t *testing.T) {
	getMock := func(db db.IDb) service.IUserManager {
		mockUser := new(service.IUserManagerMock)
		mockUser.On("UpdateUser", 1, &model.UserRequest{
			Email: "user@domain.tld", Permissions: []string{"ADMIN", "CRT_LIST", "MOD_LIST",
				"CRT_USER", "MOD_USER"},
		}).Return(nil)
		return mockUser
	}

	updateUser([]string{"1", "--email=user@domain.tld", "--permission=ADMIN,CRT_LIST,MOD_LIST,CRT_USER,MOD_USER"},
		&db.Db{}, getMock)
}
