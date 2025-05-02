package main

import (
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/cmd/goformail/cli"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"os"
)

func main() {
	// Load Basics
	configs := util.LoadConfigs("configs.cf")
	dbObj := db.InitDB(configs)

	if len(os.Args) > 1 {
		cli.NewRouter().RouteCommand(os.Args[1:], dbObj)
		return
	}
	InitApp(configs, dbObj)
}
