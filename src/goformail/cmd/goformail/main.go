package main

import (
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/forwarding"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/rest"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
)

func main() {
	configs := util.LoadConfigs("configs.cf")
	dbObj := db.InitDB(configs)
	go forwarding.LMTPService(configs, dbObj)
	auth := service.NewAuthManager(dbObj, dbObj.GetJwtSecret())
	userManager := service.NewUserManager(dbObj)
	rest.NewController(configs, dbObj, userManager, auth).Serve()
}
