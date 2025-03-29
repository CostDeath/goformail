package main

import (
	"github.com/joho/godotenv"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/interfaces"
	mailforwarding "gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/mailForwarding"
	"log"
)

func loadConfigs() map[string]string {
	configs, err := godotenv.Read("configs.cf")
	if err != nil {
		log.Fatal(err)
	}
	return configs
}

func main() {
	configs := loadConfigs()
	go mailforwarding.LMTPService(configs)
	interfaces.ServeHttp()
}
