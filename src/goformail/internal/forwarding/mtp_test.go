package forwarding

import (
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/test"
	"log"
	"net"
	"strings"
	"sync"
	"testing"
)

func TestSendResponse(t *testing.T) {
	// MOCK LISTENER
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)
	go func() {

		isSuccessful := test.SendResponseMock(waitGroup)
		if !isSuccessful {
			t.Error("Response was not able to be received")
		}

		t.Log("Goroutine finished within testSendResponse")
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

	sendResponse("Response test", conn)

}

func TestSendSuccessfulGoodBye(t *testing.T) {
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)

	go func() {
		returnedMessages := strings.Lines(test.SendGoodbyeMock(waitGroup))

		t.Log("Goroutine has finished in TestSendSuccessfulGoodbye")

		linesReceived := 0
		for message := range returnedMessages {
			switch {
			case strings.TrimSpace(message) == "452 temporarily over quota":
				linesReceived += 1
				t.Log("452 temporarily over quota - Response passed")
			case strings.TrimSpace(message) == "250 OK (Email was successfully forwarded)":
				linesReceived += 1
				t.Log("250 OK (Email was successfully forwarded)")
			case strings.TrimSpace(message) == "221 closing connection":
				linesReceived += 1
				t.Log("221 closing connection - Response passed")
			default:
				t.Error("Unexpected response was received")
			}
		}

		if linesReceived != 3 {
			t.Errorf("Not all responses were received, only %d cases passed", linesReceived)
		}
		waitGroup.Done()
	}()

	waitGroup.Wait()
	waitGroup.Add(1)
	conn, err := net.Dial("tcp", "127.0.0.1:8025")
	if err != nil {
		log.Fatal(err)
	}

	defer func(conn net.Conn) {
		waitGroup.Wait()
		if err = conn.Close(); err != nil {
			log.Fatal(err)
		}
	}(conn)

	sendGoodbye(conn, true, []string{"QUIT"})
}

func TestSendFailGoodBye(t *testing.T) {
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)

	go func() {
		returnedMessages := strings.Lines(test.SendGoodbyeMock(waitGroup))

		t.Log("Goroutine has finished in TestSendSuccessfulGoodbye")

		passedTests := 0
		for message := range returnedMessages {
			switch {
			case strings.TrimSpace(message) == "452 temporarily over quota":
				passedTests += 1
				t.Log("452 temporarily over quota - Response passed")
			case strings.TrimSpace(message) == "250 OK (Email was successfully forwarded)":
				passedTests += 1
				t.Log("250 OK (However, email was not forwarded)")
			case strings.TrimSpace(message) == "221 closing connection":
				passedTests += 1
				t.Log("221 closing connection - Response passed")
			}
		}

		if passedTests != 3 {
			t.Errorf("Not all responses were received, only %d cases passed", passedTests)
		}
		waitGroup.Done()
	}()

	waitGroup.Wait()
	waitGroup.Add(1)
	conn, err := net.Dial("tcp", "127.0.0.1:8025")
	if err != nil {
		log.Fatal(err)
	}

	defer func(conn net.Conn) {
		waitGroup.Wait()
		if err = conn.Close(); err != nil {
			log.Fatal(err)
		}
	}(conn)

	sendGoodbye(conn, true, []string{"QUIT"})
}

func TestSuccessfulMailReceiver(t *testing.T) {
	tcpSocket, err := createLMTPSocket("8024")
	if err != nil {
		log.Fatal(err)
	}
	configs := test.Configs
	waitGroup := new(sync.WaitGroup)

	defer func(tcpSocket net.Listener) {
		err = tcpSocket.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(tcpSocket)

	// MOCK MTA
	go func() {

		result := test.MailReceiverMock("LHLO example.domain", waitGroup)
		if result != "successfully received the email" {
			t.Error(result)
		}

		t.Log("Goroutine finished within TestMailReceiver")
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

func TestFailedMailReceiver(t *testing.T) {
	tcpSocket, err := createLMTPSocket("8024")
	if err != nil {
		log.Fatal(err)
	}
	configs := test.Configs
	waitGroup := new(sync.WaitGroup)

	defer func(tcpSocket net.Listener) {
		err = tcpSocket.Close()
		if err != nil {
			log.Fatal(err)
		}
		waitGroup.Wait()
	}(tcpSocket)

	// MOCK MTA
	go func() {
		result := test.MailReceiverMock("Invalid response", waitGroup)
		if result != "receiving the email was unsuccessful" {
			t.Error(result)
		}

		t.Log("Goroutine finished within TestMailReceiver")
	}()
	waitGroup.Add(1)

	conn, err := tcpSocket.Accept()
	if err != nil {
		log.Fatal(err)
	}

	_, mailError := mailReceiver(conn, 4096, configs)

	waitGroup.Wait()

	if mailError == nil {
		t.Error("An error type was not returned")
	}
}

func TestMailSenderSuccess(t *testing.T) {
	configs := test.Configs
	configs["POSTFIX_PORT"] = "8025"
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)

	// MOCK SMTP SERVER
	go func() {
		result := test.MailSenderMock("250-example.domain\n250-PIPELINING\n250 CHUNKING\n", waitGroup)
		if result != "Exited the connection" {
			t.Error(result)
		}
		t.Log("Goroutine finished within TestMailSender")
	}()
	waitGroup.Wait()
	waitGroup.Add(1)

	hasSent := mailSender("example@example.domain", "hello\n.\nQUIT\n", 4096, configs)

	if !hasSent {
		t.Error("The expected response result was not given")
	}
	waitGroup.Wait()
}

func TestMailSenderFailure(t *testing.T) {
	configs := test.Configs
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)

	// MOCK SMTP SERVER
	go func() {
		result := test.MailSenderMock("Invalid reply", waitGroup)
		if result != "Exited the connection" {
			t.Error(result)
		}
		t.Log("Goroutine finished within TestMailSender")
	}()
	waitGroup.Wait()
	waitGroup.Add(1)

	hasSent := mailSender("example@example.domain", "hello\n.\nQUIT\n", 4096, configs)

	if hasSent {
		t.Error("Function believes the email should have went through")
	}
	waitGroup.Wait()
}
