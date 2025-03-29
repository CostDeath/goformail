package forwarding

import (
	"errors"
	config "gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/configs"
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

func TestValidEmail(t *testing.T) {
	positive := "this-works@example.com"
	negative := "-this-doesnt@example.com"
	if matches := validEmail(positive); !matches {
		t.Errorf("Email is not valid when it should be: %s", positive)
	}
	if matches := validEmail(negative); matches {
		t.Errorf("Email is valid when it should not be: %s", negative)
	}
}

func TestConnectToLMTP(t *testing.T) {
	tcpSocket := createLMTPSocket("8024")
	if tcpSocket == nil {
		t.Errorf("tcpSocket is nil")
		return
	}

	defer func(tcpSocket net.Listener) {
		err := tcpSocket.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(tcpSocket)
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

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	if conn == nil {
		t.Error("Could not connect to SMTP port")
		return
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

	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

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
	tcpSocket := createLMTPSocket("8024")
	configs := config.LoadConfigs()
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
					if _, err = conn.Write([]byte("hello\n.\nmultiple full stops\n.\nQUIT\n")); err != nil {
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

	data, err := mailReceiver(conn, 4096, configs)

	waitGroup.Wait()

	if err != nil {
		t.Errorf("Error was thrown: %s", err.Error())
	} else if data.data != "hello\n.\nmultiple full stops\n.\n" {
		t.Errorf("The email data was received incorrectly\nExpected:\n hello\n.\nmultiple full stops\n.\n Received: %s\n", data.data)
	}

}

func TestMailSender(t *testing.T) {
	configs := config.LoadConfigs()
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

	hasSent := mailSender("example@example.domain", "hello", 4096, configs)

	if !hasSent {
		t.Error("The expected response result was not given")
	}
}
