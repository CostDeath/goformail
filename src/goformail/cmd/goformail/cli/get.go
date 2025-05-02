package cli

import (
	"flag"
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"log"
)

var unknownGetCommand = "Unknown command. Try 'goformail help get' for help!"

func (r *Router) RouteGetCommand(args []string, dbObj *db.Db) {
	switch args[0] {
	case "emails":
		r.getEmails(args[1:], dbObj)
	case "list":
		r.getList(args[1:], dbObj, getListManager)
	case "lists":
		r.getLists(args[1:], dbObj, getListManager)
	case "user":
		r.getUser(args[1:], dbObj, getUserManager)
	case "users":
		r.getUsers(args[1:], dbObj, getUserManager)
	default:
		fmt.Println(unknownGetCommand)
	}
}

func getEmails(args []string, dbObj db.IDb) {
	cmd := flag.NewFlagSet("emails", flag.ExitOnError)
	offset := cmd.Int("offset", 0, "offset for emails")
	list := cmd.Int("list", 0, "filter by list id")
	exhausted := cmd.Bool("exhausted", false, "show only exhausted emails")
	archived := cmd.Bool("archived", false, "show only sent emails")
	pending := cmd.Bool("pending", false, "show only emails pending approval")
	parseArgs(cmd, args)

	if cmd.NArg() != 0 {
		fmt.Println(unknownGetCommand)
		return
	}

	if *list == 0 {
		list = nil
	}

	reqs := model.EmailReqs{
		Offset:          *offset,
		List:            list,
		Archived:        *archived,
		Exhausted:       *exhausted,
		PendingApproval: *pending,
	}

	resp, err := dbObj.GetAllEmails(&reqs)
	if err != nil {
		log.Fatal(err.Err)
	}

	printObject(resp)
}

func getList(args []string, dbObj db.IDb, newManager func(db.IDb) service.IListManager) {
	cmd := flag.NewFlagSet("list", flag.ExitOnError)
	parseArgs(cmd, args)
	if cmd.NArg() != 1 {
		fmt.Println(unknownGetCommand)
		return
	}

	id := convertId(args[0])

	list, err := newManager(dbObj).GetList(id)
	if err != nil {
		log.Fatal(err.Message)
	}

	printObject(list)
}

func getLists(args []string, dbObj db.IDb, newManager func(db.IDb) service.IListManager) {
	cmd := flag.NewFlagSet("list", flag.ExitOnError)
	parseArgs(cmd, args)
	if cmd.NArg() != 0 {
		fmt.Println(unknownGetCommand)
		return
	}

	list, err := newManager(dbObj).GetAllLists()
	if err != nil {
		log.Fatal(err.Message)
	}

	printObject(list)
}

func getUser(args []string, dbObj db.IDb, newManager func(db.IDb) service.IUserManager) {
	cmd := flag.NewFlagSet("user", flag.ExitOnError)
	parseArgs(cmd, args)
	if cmd.NArg() != 1 {
		fmt.Println(unknownGetCommand)
		return
	}

	id := convertId(args[0])

	list, err := newManager(dbObj).GetUser(id)
	if err != nil {
		log.Fatal(err.Message)
	}

	printObject(list)
}

func getUsers(args []string, dbObj db.IDb, newManager func(db.IDb) service.IUserManager) {
	cmd := flag.NewFlagSet("user", flag.ExitOnError)
	parseArgs(cmd, args)
	if cmd.NArg() != 0 {
		fmt.Println(unknownGetCommand)
		return
	}

	list, err := newManager(dbObj).GetAllUsers()
	if err != nil {
		log.Fatal(err.Message)
	}

	printObject(list)
}
