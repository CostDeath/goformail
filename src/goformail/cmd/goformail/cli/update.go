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

var unknownUpdateCommand = "Unknown command. Try 'goformail help update' for help!"

func (r *Router) RouteUpdateCommand(args []string, dbObj *db.Db) {
	switch args[0] {
	case "list":
		r.updateList(args[1:], dbObj, getListManager)
	case "user":
		r.updateUser(args[1:], dbObj, getUserManager)
	default:
		fmt.Println(unknownUpdateCommand)
	}
}

func updateList(args []string, dbObj db.IDb, newManager func(db.IDb) service.IListManager) {
	if len(args) < 1 {
		fmt.Println(unknownUpdateCommand)
		return
	}
	args = append(args[1:], args[0])

	var recipients stringSlice
	var mods intSlice
	var approvedSenders stringSlice
	cmd := flag.NewFlagSet("list", flag.ExitOnError)
	name := cmd.String("name", "", "name of mailing list (without the domain)")
	cmd.Var(&recipients, "recipient", "recipients for the mailing list (can be repeated or comma-separated)")
	cmd.Var(&mods, "mod", "users allowed to update the mailing list (user ids, can be repeated or comma-separated)")
	cmd.Var(&approvedSenders, "approved",
		"emails that bypass approval for the mailing list (can be repeated or comma-separated)")
	locked := cmd.String("locked", "", "whether emails need approval (true or false)")
	parseArgs(cmd, args)
	if cmd.NArg() != 1 {
		fmt.Println(unknownUpdateCommand)
		return
	}

	if len(recipients) == 1 && recipients[0] == "" {
		recipients = stringSlice{}
	}
	if len(mods) == 1 && mods[0] == 0 {
		mods = intSlice{}
	}
	if len(approvedSenders) == 1 && approvedSenders[0] == "" {
		approvedSenders = stringSlice{}
	}

	lockedSet := false
	lockedBool := false
	if *locked != "" {
		value, e := strconv.ParseBool(*locked)
		if e != nil {
			log.Fatal("Invalid lock state. Valid states: True, False")
		}
		lockedSet = true
		lockedBool = value
	}

	id := convertId(args[len(args)-1])

	err := newManager(dbObj).UpdateList(id, &model.ListRequest{
		Name:            *name,
		Recipients:      recipients,
		Mods:            mods,
		ApprovedSenders: approvedSenders,
		Locked:          lockedBool,
	}, lockedSet)
	if err != nil {
		log.Fatal(err.Message)
	}

	fmt.Println("Successfully updated list!")
}

func updateUser(args []string, dbObj db.IDb, newManager func(db.IDb) service.IUserManager) {
	if len(args) < 1 {
		fmt.Println(unknownUpdateCommand)
		return
	}
	args = append(args[1:], args[0])

	var perms stringSlice
	cmd := flag.NewFlagSet("user", flag.ExitOnError)
	email := cmd.String("email", "", "email of user (with the domain)")
	password := cmd.String("password", "", "password for the user")
	cmd.Var(&perms, "permission",
		"permissions for user (ADMIN/CRT_USER/MOD_USER/CRT_LIST/MOD_LIST) (can be repeated or comma-separated)")
	parseArgs(cmd, args)
	if cmd.NArg() != 1 {
		fmt.Println(unknownUpdateCommand)
		return
	}

	id := convertId(args[len(args)-1])

	if len(perms) == 1 && perms[0] == "" {
		perms = stringSlice{}
	}

	err := newManager(dbObj).UpdateUser(id, &model.UserRequest{
		Email:       *email,
		Password:    *password,
		Permissions: perms,
	})
	if err != nil {
		log.Fatal(err.Message)
	}

	fmt.Println("Successfully updated list!")
}
