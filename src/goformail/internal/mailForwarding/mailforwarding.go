package forwarding

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type EmailData struct {
	rcpt          []string
	from          string
	data          string
	remainingAcks []string
}

type emailCollectionError struct {
	errorType string
	Err       error
}

func (e *emailCollectionError) Error() string {
	return fmt.Sprintf("%s: %s", e.errorType, e.Err.Error())
}

func LMTPService(configs map[string]string) {
	// TODO: Handle permissions to be able to send from mailing lists
	// need db for this
	// for now, assume all email addresses are currently valid
	fmt.Println(getCurrentTime() + " Starting LMTP Service...")
	lmtpPort := configs["LMTP_PORT"]

	tcpSocket := createLMTPSocket(lmtpPort)

	bufferSize, err := strconv.Atoi(configs["BUFFER_SIZE"])
	if err != nil {
		log.Fatal(err)
	}

	var conn net.Conn
	var mailForwardSuccess bool
	var data EmailData
	for {
		conn, err = tcpSocket.Accept()
		if err != nil {
			log.Fatal(err)
		}
		mailForwardSuccess = false

		// MAIL RECEIVER LOGIC
		data, err = mailReceiver(conn, bufferSize, configs)
		if err != nil {
			fmt.Printf("%s %s\n", getCurrentTime(), err)
			if err = conn.Close(); err != nil {
				log.Fatal(err)
			}
			continue // want to go back to start of loop
		}

		// SEND MAIL LOGIC
		if data.data != "" {
			for _, mailingList := range data.rcpt {
				mailForwardSuccess = mailSender(mailingList, data.data, bufferSize, configs)
			}
		}

		// GOODBYE ACKNOWLEDGEMENT TO RESTART
		sendGoodbye(conn, mailForwardSuccess, configs["REMAINING_ACK"])
	}
}

func getCurrentTime() string {
	t := time.Now()
	return fmt.Sprintf("[%d-%02d-%02d %02d:%02d:%02d]", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func createLMTPSocket(lmtpPort string) net.Listener {
	tcpSocket, err := net.Listen("tcp", fmt.Sprintf(":%s", lmtpPort))
	if err != nil {
		log.Fatal(err)
	}
	return tcpSocket
}

func connectToSMTP(smtpAddress string, smtpPort string) net.Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", smtpAddress, smtpPort))
	if err != nil {
		log.Fatal(getCurrentTime() + err.Error())
	}

	return conn
}

func sendResponse(resp string, conn net.Conn) {
	if _, err := conn.Write([]byte(resp)); err != nil {
		log.Fatal(err)
	}
	fmt.Println(getCurrentTime() + " S: " + resp)
}

func sendGoodbye(conn net.Conn, mailForwardSuccess bool, remainingAcks string) {
	if mailForwardSuccess {
		sendResponse("250 OK (Email was successfully forwarded)\n452 temporarily over quota\n", conn)
	} else {
		sendResponse("250 OK (However, email was not forwarded)\n452 temporarily over quota\n", conn)
	}

	messages := strings.Lines(remainingAcks)
	for message := range messages {
		if strings.TrimSpace(message) == "QUIT" {
			sendResponse("221 closing connection\n", conn)
			fmt.Println(getCurrentTime() + " S: Email successfully received, listening for more emails...")
		} else {
			if err := conn.Close(); err != nil {
				fmt.Println(getCurrentTime() + " ERROR: Unexpected response, closing connection...")
			}
		}
	}
}
