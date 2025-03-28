package mailforwarding

import (
	"github.com/joho/godotenv"
	"log"
	"regexp"
	"testing"
)

func TestGetCurrentTime(t *testing.T) {
	formattedTime := getCurrentTime()

	// EXPECTED FORMAT [YYYY-MM-DD HH:MM:SS]
	matched, err := regexp.Match(`\[([0-9]{4})(-([0-9]{2})){2}\s([0-5][0-9]:){2}[0-5][0-9]]`, []byte(formattedTime))
	if err != nil {
		log.Fatal(err)
	}
	if !matched {
		t.Errorf("Not the expected format of [YYYY-MM-DD HH:mm:ss], got %s", formattedTime)
	}
}

func TestLoadConfigs(t *testing.T) {
	configs, err := godotenv.Read("configs.cf")
	if err != nil {
		t.Error("Cannot locate config file")
	}

	//update this test if config names get modified/added
	configList := []string{"EMAIL_DOMAIN", "DEBUG_MODE", "POSTFIX_ADDRESS", "POSTFIX_PORT", "LMTP_PORT", "TIMEOUT_DURATION", "BUFFER_SIZE"}
	for _, configParam := range configList {
		if _, exists := configs[configParam]; !exists {
			t.Errorf("Config %s is not a parameter within the configs", configParam)
		}
	}
}
