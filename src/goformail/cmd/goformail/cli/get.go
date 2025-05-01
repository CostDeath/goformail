package cli

import (
	"flag"
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"log"
)

var unknownGetCommand = "Unknown command. Try 'help get' for help!"

func (r *Router) RouteGetCommand(args []string, dbObj *db.Db) {
	switch args[0] {
	case "emails":
		r.getEmails(args[1:], dbObj)
	case "list":
		r.getEmails(args[1:], dbObj)
	case "lists":
		r.getEmails(args[1:], dbObj)
	case "user":
		r.getEmails(args[1:], dbObj)
	case "users":
		r.getEmails(args[1:], dbObj)
	default:
		fmt.Println(unknownGetCommand)
	}
}

func getEmails(args []string, dbObj *db.Db) {
	if len(args) != 0 {
		fmt.Println(unknownGetCommand)
		return
	}

	offset := flag.Int("offset", 0, "offset of emails")
	list := flag.Int("list", 0, "offset of emails")
	exhausted := flag.Bool("exhausted", false, "")
	archived := flag.Bool("archived", false, "")
	pending := flag.Bool("pendingApproval", false, "")
	flag.Parse()

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

func getList(args []string, dbObj *db.Db) {
	if len(args) != 1 {
		fmt.Println(unknownGetCommand)
		return
	}

	id := convertId(args[0])

	list, err := service.NewListManager(dbObj).GetList(id)
	if err != nil {
		log.Fatal(err.Message)
	}

	printObject(list)
}

func getLists(args []string, dbObj *db.Db) {
	if len(args) != 0 {
		fmt.Println(unknownGetCommand)
		return
	}

	list, err := service.NewListManager(dbObj).GetAllLists()
	if err != nil {
		log.Fatal(err.Message)
	}

	printObject(list)
}
