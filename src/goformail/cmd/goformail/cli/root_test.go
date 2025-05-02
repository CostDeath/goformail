package cli

import (
	"github.com/stretchr/testify/assert"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"testing"
)

type mock struct {
	t      *testing.T
	called bool
	args   []string
}

func (mock *mock) mockFun(args []string, dbObj db.IDb) {
	mock.called = true
	mock.args = args
	assert.Equal(mock.t, &db.Db{}, dbObj)
}

func (mock *mock) mockListFun(args []string, dbObj db.IDb, newManager func(db.IDb) service.IListManager) {
	mock.called = true
	mock.args = args
	assert.Equal(mock.t, &db.Db{}, dbObj)
	assert.Equal(mock.t, service.NewListManager(dbObj), newManager(dbObj))
}

func (mock *mock) mockUserFun(args []string, dbObj db.IDb, newManager func(db.IDb) service.IUserManager) {
	mock.called = true
	mock.args = args
	assert.Equal(mock.t, &db.Db{}, dbObj)
	assert.Equal(mock.t, service.NewUserManager(dbObj), newManager(dbObj))
}

func TestControllerGetEmails(t *testing.T) {
	args := []string{"get", "emails", "--exhausted", "--list=1"}

	mock := &mock{t: t}
	router := Router{getEmails: mock.mockFun}
	router.RouteCommand(args, &db.Db{})
	assert.True(t, mock.called)
	assert.Equal(t, args[2:], mock.args)
}

func TestControllerGetList(t *testing.T) {
	args := []string{"get", "list", "1"}

	mock := &mock{t: t}
	router := Router{getList: mock.mockListFun}
	router.RouteCommand(args, &db.Db{})
	assert.True(t, mock.called)
	assert.Equal(t, args[2:], mock.args)
}

func TestControllerGetLists(t *testing.T) {
	args := []string{"get", "lists"}

	mock := &mock{t: t}
	router := Router{getLists: mock.mockListFun}
	router.RouteCommand(args, &db.Db{})
	assert.True(t, mock.called)
	assert.Equal(t, args[2:], mock.args)
}

func TestControllerGetUser(t *testing.T) {
	args := []string{"get", "user", "1"}

	mock := &mock{t: t}
	router := Router{getUser: mock.mockUserFun}
	router.RouteCommand(args, &db.Db{})
	assert.True(t, mock.called)
	assert.Equal(t, args[2:], mock.args)
}

func TestControllerGetUsers(t *testing.T) {
	args := []string{"get", "users"}

	mock := &mock{t: t}
	router := Router{getUsers: mock.mockUserFun}
	router.RouteCommand(args, &db.Db{})
	assert.True(t, mock.called)
	assert.Equal(t, args[2:], mock.args)
}

func TestControllerCreateList(t *testing.T) {
	args := []string{"create", "list", "--name=name"}

	mock := &mock{t: t}
	router := Router{createList: mock.mockListFun}
	router.RouteCommand(args, &db.Db{})
	assert.True(t, mock.called)
	assert.Equal(t, args[2:], mock.args)
}

func TestControllerCreateUser(t *testing.T) {
	args := []string{"create", "user", "--email=example@domain.tld"}

	mock := &mock{t: t}
	router := Router{createUser: mock.mockUserFun}
	router.RouteCommand(args, &db.Db{})
	assert.True(t, mock.called)
	assert.Equal(t, args[2:], mock.args)
}

func TestControllerUpdateList(t *testing.T) {
	args := []string{"update", "list", "1", "--name=name"}

	mock := &mock{t: t}
	router := Router{updateList: mock.mockListFun}
	router.RouteCommand(args, &db.Db{})
	assert.True(t, mock.called)
	assert.Equal(t, args[2:], mock.args)
}

func TestControllerUpdateUser(t *testing.T) {
	args := []string{"update", "user", "1", "--email=example@domain.tld"}

	mock := &mock{t: t}
	router := Router{updateUser: mock.mockUserFun}
	router.RouteCommand(args, &db.Db{})
	assert.True(t, mock.called)
	assert.Equal(t, args[2:], mock.args)
}

func TestControllerDeleteList(t *testing.T) {
	args := []string{"delete", "list", "1"}

	mock := &mock{t: t}
	router := Router{deleteList: mock.mockListFun}
	router.RouteCommand(args, &db.Db{})
	assert.True(t, mock.called)
	assert.Equal(t, args[2:], mock.args)
}

func TestControllerDeleteUser(t *testing.T) {
	args := []string{"delete", "user", "1"}

	mock := &mock{t: t}
	router := Router{deleteUser: mock.mockUserFun}
	router.RouteCommand(args, &db.Db{})
	assert.True(t, mock.called)
	assert.Equal(t, args[2:], mock.args)
}

func TestControllerApproveEmail(t *testing.T) {
	args := []string{"approve", "email", "1"}

	mock := &mock{t: t}
	router := Router{approveEmail: mock.mockFun}
	router.RouteCommand(args, &db.Db{})
	assert.True(t, mock.called)
	assert.Equal(t, args[2:], mock.args)
}
