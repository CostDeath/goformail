package mailforwarding

import (
	"github.com/joho/godotenv"
	"log"
	"net"
	"regexp"
	"sync"
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

func TestConnectToLMTP(t *testing.T) {
	tcpSocket := connectToLMTP("8024")
	if tcpSocket == nil {
		t.Errorf("tcpSocket is nil")
		return
	}

	if err := tcpSocket.Close(); err != nil {
		t.Error(err)
	}
}

func TestConnectToSMTP(t *testing.T) {
	waitGroup := new(sync.WaitGroup)
	// MOCK Listener
	go func() {
		tcpSocket, err := net.Listen("tcp", "127.0.0.1:8025")
		if err != nil {
			log.Fatal(err)
		}
		defer func(tcpSocket net.Listener) {
			err = tcpSocket.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(tcpSocket)

		waitGroup.Done()
		_, err = tcpSocket.Accept()
	}()

	waitGroup.Add(1)
	waitGroup.Wait()
	conn := connectToSMTP("127.0.0.1", "8025")

	if conn == nil {
		t.Error("Could not connect to SMTP port")
		return
	}

	if err := conn.Close(); err != nil {
		t.Error(err)
	}
}

func TestSendResponse(t *testing.T) {
	// MOCK LISTENER
	waitGroup := new(sync.WaitGroup)
	go func() {
		tcpSocket, err := net.Listen("tcp", "127.0.0.1:8025")
		if err != nil {
			log.Fatal(err)
		}
		defer func(tcpSocket net.Listener) {
			err = tcpSocket.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(tcpSocket)

		waitGroup.Done()
		conn, err := tcpSocket.Accept()
		if err != nil {
			log.Fatal(err)
		}

		sendResponse("Response test", conn)
	}()

	waitGroup.Add(1)
	waitGroup.Wait()
	conn, err := net.Dial("tcp", "127.0.0.1:8025")
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 1024)

	size, err := conn.Read(buffer)
	if err != nil {
		t.Error(err)
	}

	if string(buffer[:size]) != "Response test" {
		t.Error("There were no responses/response was wrong")
	}
}
