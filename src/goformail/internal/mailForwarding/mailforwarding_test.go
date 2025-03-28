package mailforwarding

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
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

func TestSendGoodBye(t *testing.T) {
	waitGroup := new(sync.WaitGroup)
	go func() {
		tcpSocket, err := net.Listen("tcp", "127.0.0.1:8025")
		if err != nil {
			log.Fatal(err)
		}

		waitGroup.Done()
		conn, err := tcpSocket.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// POSITIVE CASE
		sendGoodbye(conn, true, "QUIT")

		// NEGATIVE CASE
		sendGoodbye(conn, false, "QUIT")
	}()
	waitGroup.Add(1)
	waitGroup.Wait()

	conn, err := net.Dial("tcp", "127.0.0.1:8025")
	if err != nil {
		log.Fatal(err)
	}

	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {

		}
	}(conn)

	buffer := make([]byte, 4096)
	var size int
	var collected string

	for i := 0; i < 4; i++ {
		if err = conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			log.Fatal(err)
		}

		size, err = conn.Read(buffer)
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			t.Error("Connection timed out, messages were not sent through connection")
			return
		}

		collected += string(buffer[:size])
	}

	messages := strings.Lines(collected)
	passedTests := 0
	for message := range messages {
		switch {
		case strings.TrimSpace(message) == "250 OK (Email was successfully forwarded)":
			passedTests += 1
		case strings.TrimSpace(message) == "452 temporarily over quota":
			passedTests += 1
		case strings.TrimSpace(message) == "250 OK (However, email was not forwarded)":
			passedTests += 1
		case strings.TrimSpace(message) == "221 closing connection":
			passedTests += 1
		}
	}
	if passedTests != 6 {
		t.Errorf("Not all cases were resolved for sending goodbye acknowledgements (only %d passed)", passedTests)
	}

}
