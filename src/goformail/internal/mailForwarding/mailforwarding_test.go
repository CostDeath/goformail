package forwarding

import (
	"errors"
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
	waitGroup.Add(1)
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
			waitGroup.Done()
		}(tcpSocket)

		waitGroup.Done()
		_, err = tcpSocket.Accept()
		if err != nil {
			log.Fatal(err)
		}
	}()

	waitGroup.Wait()
	waitGroup.Add(1)
	conn := connectToSMTP("127.0.0.1", "8025")

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
		waitGroup.Wait()
	}(conn)
}

func TestSendResponse(t *testing.T) {
	// MOCK LISTENER
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)
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
			waitGroup.Done() // ensure go routine finishes first
		}(tcpSocket)

		waitGroup.Done()
		conn, err := tcpSocket.Accept()
		if err != nil {
			log.Fatal(err)
		}

		sendResponse("Response test", conn)
	}()

	waitGroup.Wait()
	waitGroup.Add(1)
	conn, err := net.Dial("tcp", "127.0.0.1:8025")
	if err != nil {
		log.Fatal(err)
	}

	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
		waitGroup.Wait() // Ensure go routine function finishes first
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
	waitGroup.Add(1)
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
			waitGroup.Done()
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
	waitGroup.Wait()
	waitGroup.Add(1)

	conn, err := net.Dial("tcp", "127.0.0.1:8025")
	if err != nil {
		log.Fatal(err)
	}

	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
		waitGroup.Wait()
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
