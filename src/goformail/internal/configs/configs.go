package config

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadConfigs() map[string]string {
	configs, err := godotenv.Read("configs.cf")
	if err != nil {
		log.Fatal(err)
	}
	return configs
}
