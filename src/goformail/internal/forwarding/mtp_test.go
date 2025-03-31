package forwarding

import (
	"errors"
	config "gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal"
	"log"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

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
			t.Log("Goroutine finished within testSendResponse")
			waitGroup.Done() // ensure go routine finishes first
		}(tcpSocket)

		waitGroup.Done()
		conn, err := tcpSocket.Accept()
		if err != nil {
			log.Fatal(err)
		}

		sendResponse("Response test", conn)
		buffer := make([]byte, 200)
		if _, err = conn.Read(buffer); err != nil {
			log.Fatal(err)
		}
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

	if _, err = conn.Write([]byte("Exit")); err != nil {
		log.Fatal("Could not write to connection")
	}

	if string(buffer[:size]) != "Response test" {
		t.Error("There were no responses/response was wrong")
	}
}

func TestMailReceiver(t *testing.T) {
	tcpSocket := createLMTPSocket("8024")
	configs := config.LoadConfigs("../../configs.cf")
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
			t.Log("Goroutine finished within TestMailReceiver")
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
	configs := config.LoadConfigs("../../configs.cf")
	configs["POSTFIX_PORT"] = "8025"
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)

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
			t.Log("Goroutine finished within TestMailSender")
			waitGroup.Done()
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
				t.Log(message)
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
	waitGroup.Wait()
	waitGroup.Add(1)

	hasSent := mailSender("example@example.domain", "hello\n", 4096, configs)

	if !hasSent {
		t.Error("The expected response result was not given")
	}
	waitGroup.Wait()
}
