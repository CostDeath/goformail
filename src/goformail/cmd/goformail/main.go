package main

import (
	config "gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/configs"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/interfaces"
	forwarding "gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/mailForwarding"
)

func main() {
	configs := config.LoadConfigs()
	go forwarding.LMTPService(configs)
	interfaces.ServeHttp()
}
