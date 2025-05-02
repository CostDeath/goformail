package cli

import (
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"testing"
)

func TestGetEmails(t *testing.T) {
	mockDb := new(db.IDbMock)
	mockDb.On("GetAllEmails", &model.EmailReqs{}).Return(&model.EmailResponse{})

	getEmails([]string{}, mockDb)
}

func TestGetEmailsAllSet(t *testing.T) {
	num := 1
	mockDb := new(db.IDbMock)
	mockDb.On("GetAllEmails", &model.EmailReqs{
		Offset: 10, List: &num, Exhausted: true, PendingApproval: true, Archived: true,
	}).Return(&model.EmailResponse{})

	getEmails([]string{"--offset=10", "--list=1", "--exhausted", "--archived", "--pending"}, mockDb)
}

func TestGetList(t *testing.T) {
	getMock := func(db db.IDb) service.IListManager {
		mockList := new(service.IListManagerMock)
		mockList.On("GetList", 1).Return(&model.ListResponse{})
		return mockList
	}

	getList([]string{"1"}, &db.Db{}, getMock)
}

func TestGetLists(t *testing.T) {
	getMock := func(db db.IDb) service.IListManager {
		mockList := new(service.IListManagerMock)
		mockList.On("GetAllLists").Return(&[]*model.ListResponse{})
		return mockList
	}

	getLists([]string{}, &db.Db{}, getMock)
}

func TestGetUser(t *testing.T) {
	getMock := func(db db.IDb) service.IUserManager {
		mockUser := new(service.IUserManagerMock)
		mockUser.On("GetUser", 1).Return(&model.UserResponse{})
		return mockUser
	}

	getUser([]string{"1"}, &db.Db{}, getMock)
}

func TestGetUsers(t *testing.T) {
	getMock := func(db db.IDb) service.IUserManager {
		mockUser := new(service.IUserManagerMock)
		mockUser.On("GetAllUsers", 1).Return(&[]*model.UserResponse{})
		return mockUser
	}

	getUser([]string{}, &db.Db{}, getMock)
}
