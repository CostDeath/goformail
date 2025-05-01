package main

import (
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/mail"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/rest"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
)

func main() {
	// Basics
	configs := util.LoadConfigs("configs.cf")
	dbObj := db.InitDB(configs)

	// Services
	list := service.NewListManager(dbObj)
	user := service.NewUserManager(dbObj)
	auth := service.NewAuthManager(dbObj, dbObj.GetJwtSecret())

	// Mail
	mtp := mail.NewMtpHandler(configs)
	sender := mail.NewEmailSender(mtp, dbObj, configs)

	// Start
	go sender.Loop()
	go mail.NewEmailReceiver(mtp, sender, configs).Loop()
	rest.NewController(list, user, auth).Serve(configs["HTTP_PORT"])
}
