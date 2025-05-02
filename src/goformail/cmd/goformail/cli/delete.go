package cli

import (
	"flag"
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"log"
)

var unknownDeleteCommand = "Unknown command. Try 'goformail help delete' for help!"

func (r *Router) RouteDeleteCommand(args []string, dbObj *db.Db) {
	switch args[0] {
	case "list":
		r.deleteList(args[1:], dbObj, getListManager)
	case "user":
		r.deleteUser(args[1:], dbObj, getUserManager)
	default:
		fmt.Println(unknownDeleteCommand)
	}
}

func deleteList(args []string, dbObj db.IDb, newManager func(db.IDb) service.IListManager) {
	cmd := flag.NewFlagSet("list", flag.ExitOnError)
	parseArgs(cmd, args)
	if cmd.NArg() != 1 {
		fmt.Println(unknownDeleteCommand)
		return
	}

	id := convertId(args[0])

	err := newManager(dbObj).DeleteList(id)
	if err != nil {
		log.Fatal(err.Message)
	}

	fmt.Println("Successfully deleted list!")
}

func deleteUser(args []string, dbObj db.IDb, newManager func(db.IDb) service.IUserManager) {
	cmd := flag.NewFlagSet("user", flag.ExitOnError)
	parseArgs(cmd, args)
	if cmd.NArg() != 1 {
		fmt.Println(unknownDeleteCommand)
		return
	}

	id := convertId(args[0])

	err := newManager(dbObj).DeleteUser(id)
	if err != nil {
		log.Fatal(err.Message)
	}

	fmt.Println("Successfully deleted user!")
}
