package util

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadConfigs(filePath string) map[string]string {
	configs, err := godotenv.Read(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return configs
}
