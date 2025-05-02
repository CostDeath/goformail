package cli

import (
	"flag"
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"log"
	"strconv"
)

var unknownCreateCommand = "Unknown command. Try 'goformail help create' for help!"

func (r *Router) RouteCreateCommand(args []string, dbObj *db.Db) {
	switch args[0] {
	case "list":
		r.createList(args[1:], dbObj, getListManager)
	case "user":
		r.createUser(args[1:], dbObj, getUserManager)
	default:
		fmt.Println(unknownCreateCommand)
	}
}

func createList(args []string, dbObj db.IDb, newManager func(db.IDb) service.IListManager) {
	recipients := stringSlice{}
	mods := intSlice{}
	approvedSenders := stringSlice{}
	cmd := flag.NewFlagSet("list", flag.ExitOnError)
	name := cmd.String("name", "", "name of mailing list (without the domain)")
	cmd.Var(&recipients, "recipient", "recipients for the mailing list (can be repeated or comma-separated)")
	cmd.Var(&mods, "mod", "users allowed to update the mailing list (user ids, can be repeated or comma-separated)")
	cmd.Var(&approvedSenders, "approved",
		"emails that bypass approval for the mailing list (can be repeated or comma-separated)")
	locked := cmd.String("locked", "false", "whether emails need approval (true or false)")
	parseArgs(cmd, args)
	if cmd.NArg() != 0 {
		fmt.Println(unknownCreateCommand)
		return
	}

	lockedBool, e := strconv.ParseBool(*locked)
	if e != nil {
		log.Fatal("Invalid lock state. Valid states: True, False")
	}

	id, err := newManager(dbObj).CreateList(&model.ListRequest{
		Name:            *name,
		Recipients:      recipients,
		Mods:            mods,
		ApprovedSenders: approvedSenders,
		Locked:          lockedBool,
	})
	if err != nil {
		log.Fatal(err.Message)
	}

	fmt.Printf("Successfully created list! Id: %d\n", id)
}

func createUser(args []string, dbObj db.IDb, newManager func(db.IDb) service.IUserManager) {
	perms := stringSlice{}
	cmd := flag.NewFlagSet("user", flag.ExitOnError)
	email := cmd.String("email", "", "email of user (with the domain)")
	password := cmd.String("password", "", "password for the user")
	cmd.Var(&perms, "permission",
		"permissions for user (ADMIN/CRT_USER/MOD_USER/CRT_LIST/MOD_LIST) (can be repeated or comma-separated)")
	parseArgs(cmd, args)
	if cmd.NArg() != 0 {
		fmt.Println(unknownCreateCommand)
		return
	}

	id, err := newManager(dbObj).CreateUser(&model.UserRequest{
		Email:       *email,
		Password:    *password,
		Permissions: perms,
	})
	if err != nil {
		log.Fatal(err.Message)
	}

	fmt.Printf("Successfully created user! Id: %d\n", id)
}
