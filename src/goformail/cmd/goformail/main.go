package main

import (
	config "gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/forwarding"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/interfaces"
)

func main() {
	configs := config.LoadConfigs("configs.cf")
	go forwarding.LMTPService(configs)
	interfaces.ServeHttp()
}
