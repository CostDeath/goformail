package cli

import (
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"log"
)

func (r *Router) RouteApproveCommand(args []string, dbObj *db.Db) {
	switch args[0] {
	case "email":
		r.getEmails(args[1:], dbObj)
	default:
		fmt.Println(unknownGetCommand)
	}
}

func approveEmail(args []string, dbObj *db.Db) {
	if len(args) != 1 {
		fmt.Println(unknownGetCommand)
		return
	}

	id := convertId(args[0])

	err := dbObj.SetEmailAsApproved(id)
	if err != nil {
		if err.Code == db.ErrNoRows {
			log.Fatalf("Invalid id '%s'\n'", args[0])
		}
		log.Fatal("An error occurred: ", err.Err)
	}

	fmt.Println("Successfully approved list!")
}
