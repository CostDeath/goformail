package cli

import (
	"encoding/json"
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"log"
	"strconv"
)

type getEmailsFunc func([]string, *db.Db)
type approveEmailFunc func(args []string, dbObj *db.Db)

type Router struct {
	getEmails    getEmailsFunc
	approveEmail approveEmailFunc
}

func NewRouter() *Router {
	return &Router{
		getEmails:    getEmails,
		approveEmail: approveEmail,
	}
}

func (r *Router) RouteCommand(args []string, dbObj *db.Db) {
	switch args[0] {
	case "get":
		r.RouteGetCommand(args[1:], dbObj)
	case "approve":
		r.RouteApproveCommand(args[1:], dbObj)
	default:
		fmt.Println("Unknown command. Try 'goformail help' for help!")
	}
}

func convertId(id string) int {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Fatalf("Invalid id '%s'\n'", id)
	}
	return idInt
}

func printObject(object interface{}) {
	marshal, err := json.MarshalIndent(object, "", "\t")
	if err != nil {
		log.Fatal("Error printing payload:", err)
	}
	fmt.Println(string(marshal))
}
