package main

import (
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/interfaces"
	mailforwarding "gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/mailForwarding"
)

func main() {
	go mailforwarding.LMTPService()
	interfaces.ServeHttp()
}
