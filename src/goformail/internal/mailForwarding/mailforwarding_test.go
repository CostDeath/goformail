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
			log.Fatal(err)
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

func TestMailReceiver(t *testing.T) {
	tcpSocket := connectToLMTP("8024")
	configs := loadConfigs()
	waitGroup := new(sync.WaitGroup)

	defer func(tcpSocket net.Listener) {
		err := tcpSocket.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(tcpSocket)

	// MOCK MTA
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:8024")
		if err != nil {
			log.Fatal(err)
		}

		defer func(conn net.Conn) {
			err = conn.Close()
			if err != nil {
				log.Fatal(err)
			}
			waitGroup.Done()
		}(conn)

		for {
			buffer := make([]byte, 4096)
			var size int

			if err = conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
				log.Fatal(err)
			}

			size, err = conn.Read(buffer)
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				t.Error("Connection timed out, MailReceiver cannot seem to send messages through the connection")
				return
			}

			messages := strings.Lines(string(buffer[:size]))
			for message := range messages {
				switch {
				case strings.TrimSpace(message) == "220 LMTP Server Ready":
					if _, err = conn.Write([]byte("LHLO example.domain")); err != nil {
						log.Fatal(err)
					}
				case strings.TrimSpace(message) == "250 SIZE":
					if _, err = conn.Write([]byte("MAIL FROM:<testing@example.domain>\nRCPT TO:<recipient@example.domain>\nDATA\n")); err != nil {
						log.Fatal(err)
					}
				case strings.TrimSpace(message) == "354 Start mail input; end with <CRLF>.<CRLF>":
					if _, err = conn.Write([]byte("hello\n.\nQUIT\n")); err != nil {
						log.Fatal(err)
					}
					return
				}
			}

		}
	}()
	waitGroup.Add(1)

	conn, err := tcpSocket.Accept()
	if err != nil {
		log.Fatal(err)
	}

	data := MailReceiver(conn, 4096, configs)

	waitGroup.Wait()

	expectedKeys := []string{"EMAIL_DATA", "RCPTS", "REMAINING_ACK"}
	for _, key := range expectedKeys {
		if _, exists := data[key]; !exists {
			t.Errorf("%s does not exist within the data collected from MailReceiver", key)
		}
	}
}

func TestMailSender(t *testing.T) {
	configs := loadConfigs()
	configs["POSTFIX_PORT"] = "8025"
	waitGroup := new(sync.WaitGroup)

	// MOCK SMTP SERVER
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

		// Initial greeting
		if _, err = conn.Write([]byte("Example SMTP Server Greeting\n")); err != nil {
			log.Fatal(err)
		}

		for {
			buffer := make([]byte, 4096)
			var size int

			if err = conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
				log.Fatal(err)
			}

			size, err = conn.Read(buffer)
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				t.Error("Connection timed out, MailSender cannot seem to send responses through the connection")
				return
			}

			messages := strings.Lines(string(buffer[:size]))
			for message := range messages {
				switch {
				case strings.HasPrefix(message, "EHLO"):
					if _, err = conn.Write([]byte("250-example.domain\n250-PIPELINING\n250 CHUNKING\n")); err != nil {
						log.Fatal(err)
					}
				case strings.HasPrefix(message, "MAIL FROM"):
					if _, err = conn.Write([]byte("250 2.1.0 OK\n")); err != nil {
						log.Fatal(err)
					}
				case strings.HasPrefix(message, "RCPT TO"):
					// assume all email addresses are valid
					if _, err = conn.Write([]byte("250 2.1.5 OK\n")); err != nil {
						log.Fatal(err)
					}
				case strings.TrimSpace(message) == "DATA":
					if _, err = conn.Write([]byte("354 Start mail input; end with <CRLF>.<CRLF>\n")); err != nil {
						log.Fatal(err)
					}
				case strings.TrimSpace(message) == "QUIT":
					return
				}
			}
		}

	}()
	waitGroup.Add(1)
	waitGroup.Wait()

	hasSent := MailSender("example@example.domain", "hello", 4096, configs)

	if !hasSent {
		t.Error("The expected response result was not given")
	}
}
