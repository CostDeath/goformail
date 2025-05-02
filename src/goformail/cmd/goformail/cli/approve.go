package cli

import (
	"flag"
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"log"
)

var unknownApproveCommand = "Unknown command. Try 'goformail help approve' for help!"

func (r *Router) RouteApproveCommand(args []string, dbObj *db.Db) {
	switch args[0] {
	case "email":
		r.approveEmail(args[1:], dbObj)
	default:
		fmt.Println(unknownApproveCommand)
	}
}

func approveEmail(args []string, dbObj db.IDb) {
	cmd := flag.NewFlagSet("approve", flag.ExitOnError)
	parseArgs(cmd, args)
	if cmd.NArg() != 1 {
		fmt.Println(unknownApproveCommand)
		return
	}

	id := convertId(args[0])

	err := dbObj.SetEmailAsApproved(id)
	if err != nil {
		if err.Code == db.ErrNoRows {
			log.Fatalf("No email with id '%s'\n", args[0])
		}
		log.Fatal("An error occurred: ", err.Err)
	}

	fmt.Println("Successfully approved list!")
}
