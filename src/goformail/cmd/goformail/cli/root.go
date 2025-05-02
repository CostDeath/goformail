package cli

import (
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
)

type cmdFunc func([]string, db.IDb)
type listFunc func([]string, db.IDb, func(db.IDb) service.IListManager)
type userFunc func([]string, db.IDb, func(db.IDb) service.IUserManager)

type Router struct {
	getEmails    cmdFunc
	getList      listFunc
	getLists     listFunc
	getUser      userFunc
	getUsers     userFunc
	createList   listFunc
	createUser   userFunc
	updateList   listFunc
	updateUser   userFunc
	deleteList   listFunc
	deleteUser   userFunc
	approveEmail cmdFunc
}

func NewRouter() *Router {
	return &Router{
		getEmails:    getEmails,
		getList:      getList,
		getLists:     getLists,
		getUser:      getUser,
		getUsers:     getUsers,
		createList:   createList,
		createUser:   createUser,
		updateList:   updateList,
		updateUser:   updateUser,
		deleteList:   deleteList,
		deleteUser:   deleteUser,
		approveEmail: approveEmail,
	}
}

func (r *Router) RouteCommand(args []string, dbObj *db.Db) {
	switch args[0] {
	case "get":
		r.RouteGetCommand(args[1:], dbObj)
	case "create":
		r.RouteCreateCommand(args[1:], dbObj)
	case "update":
		r.RouteUpdateCommand(args[1:], dbObj)
	case "delete":
		r.RouteDeleteCommand(args[1:], dbObj)
	case "approve":
		r.RouteApproveCommand(args[1:], dbObj)
	default:
		fmt.Println("Unknown command. Try 'goformail help' for help!")
	}
}
