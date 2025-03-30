package config

import (
	"testing"
)

func TestLoadConfigs(t *testing.T) {
	configs := LoadConfigs("../configs.cf")

	//update this test if config names get modified/added
	configList := []string{"EMAIL_DOMAIN", "DEBUG_MODE", "POSTFIX_ADDRESS", "POSTFIX_PORT", "LMTP_PORT", "TIMEOUT_DURATION", "BUFFER_SIZE"}
	for _, configParam := range configList {
		if _, exists := configs[configParam]; !exists {
			t.Errorf("Config %s is not a parameter within the configs", configParam)
		}
	}
}
