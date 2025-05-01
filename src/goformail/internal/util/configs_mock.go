package util

var MockConfigs = map[string]string{
	"EMAIL_DOMAIN":     "example.domain",
	"DEBUG_MODE":       "false",
	"POSTFIX_ADDRESS":  "127.0.0.1",
	"POSTFIX_PORT":     "8025",
	"LMTP_PORT":        "8024",
	"TIMEOUT_DURATION": "5",
	"BUFFER_SIZE":      "4096",
	"ORIGINAL_SENDER":  "true",
	"CHECK_FREQUENCY":  "60",
	"SQL_ADDRESS":      "127.0.0.1",
	"SQL_PORT":         "6432",
	"SQL_USER":         "goformail",
	"SQL_PASSWORD":     "password",
	"SQL_DB_NAME":      "goformail",
	"HTTP_PORT":        "8000",
}
