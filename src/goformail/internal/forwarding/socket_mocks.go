package forwarding

import (
	"errors"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func ConnectSMTPSocketMock(waitGroup *sync.WaitGroup) {
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
}

func ConnectFailSMTPSocketMock(waitGroup *sync.WaitGroup) {
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
}

func SendResponseMock(waitGroup *sync.WaitGroup) bool {
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

	if err = conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		log.Fatal(err)
	}

	var netErr net.Error
	buffer := make([]byte, 1024)
	if _, err = conn.Read(buffer); err != nil {
		log.Fatal(err)
	} else if errors.As(err, &netErr) && netErr.Timeout() {
		return false
	}
	return true
}

func SendGoodbyeMock(waitGroup *sync.WaitGroup) string {
	tcpSocket, err := net.Listen("tcp", "127.0.0.1:8025")
	if err != nil {
		log.Fatal(err)
	}

	defer func(tcpSocket net.Listener) {
		if err = tcpSocket.Close(); err != nil {
			log.Fatal(err)
		}
	}(tcpSocket)

	waitGroup.Done()
	conn, err := tcpSocket.Accept()
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 4096)
	var size int
	var messages string
	var netErr net.Error
	for {
		if err = conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
			log.Fatal(err)
		}
		size, err = conn.Read(buffer)
		if errors.As(err, &netErr) && netErr.Timeout() {
			break
		}

		messages += string(buffer[:size])
	}
	return messages
}

func MailReceiverMock(greeting string, waitGroup *sync.WaitGroup) string {
	var conn net.Conn
	var err error
	conn, err = net.Dial("tcp", "127.0.0.1:8024")
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
			return "Connection timed out, mail receiver seems to not be receiving any responses or could not give any responses"
		}

		messages := strings.Lines(string(buffer[:size]))
		for message := range messages {
			switch {
			case strings.TrimSpace(message) == "220 LMTP Server Ready":
				if _, err = conn.Write([]byte(greeting)); err != nil {
					log.Fatal(err)
				}
			case strings.TrimSpace(message) == "500 Error: command not recognised":
				return "receiving the email was unsuccessful" // MailReceiver should return with an error
			case strings.TrimSpace(message) == "250 SIZE":
				msg := "MAIL FROM:<testing@example.domain>\nRCPT TO:<recipient@example.domain>\nDATA\n"
				if _, err = conn.Write([]byte(msg)); err != nil {
					log.Fatal(err)
				}
			case strings.TrimSpace(message) == "354 Start mail input; end with <CRLF>.<CRLF>":
				if _, err = conn.Write([]byte("hello\n.\nmultiple full stops\n.\nQUIT\n")); err != nil {
					log.Fatal(err)
				}
				return "successfully received the email"
			}
		}

	}
}

func MailSenderMock(greeting string, waitGroup *sync.WaitGroup) string {
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
			return "Connection timed out, MailSender cannot seem to send responses through the connection"
		}

		messages := strings.Lines(string(buffer[:size]))
		for message := range messages {
			switch {
			case strings.HasPrefix(message, "EHLO"):
				if _, err = conn.Write([]byte(greeting)); err != nil {
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
				return "Exited the connection"
			}
		}
	}

}
