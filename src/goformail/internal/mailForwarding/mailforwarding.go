package mailforwarding

import (
	"errors"
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
		data, err = MailReceiver(conn, bufferSize, configs)
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
				mailForwardSuccess = MailSender(mailingList, data.data, bufferSize, configs)
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

func MailReceiver(conn net.Conn, bufferSize int, configs map[string]string) (EmailData, error) {
	domainName := configs["EMAIL_DOMAIN"]
	debugMode := configs["DEBUG_MODE"]

	data := EmailData{}

	if _, err := conn.Write([]byte("220 LMTP Server Ready\n")); err != nil {
		log.Fatal(err)
	}
	fmt.Println(getCurrentTime() + "Initialising LMTP greeting")
	inData := false
	for {
		var size int
		buffer := make([]byte, bufferSize)

		size, err := conn.Read(buffer)
		if err != nil {
			return data, &emailCollectionError{"READ_ERROR", err}
		}

		messages := strings.Lines(string(buffer[:size]))

		var emailMessage string
		for message := range messages {
			if debugMode == "true" {
				fmt.Print("POSTFIX: " + message)
			}
			switch {
			case strings.HasPrefix(message, "LHLO"):
				sendResponse(fmt.Sprintf("250-%s\n250-PIPELINING\n250 SIZE\n", domainName), conn)
			case strings.HasPrefix(message, "MAIL FROM"):
				email := strings.Fields(message)[1][6:]
				email = email[:len(email)-1]
				data.from = email
				sendResponse("250 OK\n", conn)
			case strings.HasPrefix(message, "RCPT TO"):
				email := strings.Fields(message)[1][4:]
				email = email[:len(email)-1]
				data.rcpt = append(data.rcpt, email)
				sendResponse("250 OK\n", conn)
			case strings.TrimSpace(message) == "DATA":
				sendResponse("354 Start mail input; end with <CRLF>.<CRLF>\n", conn)
				inData = true
			case strings.TrimSpace(message) == "QUIT":
				data.data = emailMessage
				data.remainingAcks = append(data.remainingAcks, message)
				return data, nil
			case inData:
				emailMessage += message
			default:
				return data, &emailCollectionError{"UNEXPECTED_RESPONSE_ERROR", errors.New(message)}
			}
		}
	}

}

func MailSender(mailingList string, emailData string, bufferSize int, configs map[string]string) bool {
	addr := configs["POSTFIX_ADDRESS"]
	port := configs["POSTFIX_PORT"]
	domainName := configs["EMAIL_DOMAIN"]
	debugMode := configs["DEBUG_MODE"]
	timeoutDuration, err := time.ParseDuration(configs["TIMEOUT_DURATION"] + "s")
	if err != nil {
		fmt.Println(getCurrentTime() + " ERROR: Could not parse timeout duration: " + err.Error())
		return false
	}

	conn := connectToSMTP(addr, port)

	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			log.Fatal(getCurrentTime() + err.Error())
		}
	}(conn)

	fmt.Println(getCurrentTime() + " S: Initiating sending email acknowledgements...")
	initial := true

	for {
		var size int
		buffer := make([]byte, bufferSize)

		err = conn.SetDeadline(time.Now().Add(timeoutDuration * time.Second)) // Time out after 5 seconds

		size, err = conn.Read(buffer)

		// handle timeout
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			fmt.Println(getCurrentTime() + " S: MTA not responding, is the correct port configured?")
			return false
		}

		messages := strings.Lines(string(buffer[:size]))

		for message := range messages {
			if debugMode == "true" {
				fmt.Print("POSTFIX: " + message)
			}
			switch {
			case initial:
				sendResponse(fmt.Sprintf("EHLO %s\n", domainName), conn)

				size, err = conn.Read(buffer)
				if err != nil {
					log.Fatal(err)
				}
				sendResponse(fmt.Sprintf("MAIL FROM: %s@%s\n", mailingList, domainName), conn)

				initial = false
			case strings.HasPrefix(message, "250 2.1.0"):
				// TODO: Similar to above...
				recipients := []string{"sdk194", "dags", "nonExistent"}
				for _, recipient := range recipients {
					sendResponse(fmt.Sprintf("RCPT TO: %s@%s\n", recipient, domainName), conn)
					size, err = conn.Read(buffer)
					message = string(buffer[:size])
					if strings.HasPrefix(message, "550 5.1.1") {
						fmt.Println(getCurrentTime() + message)
					}
				}
				sendResponse("DATA\n", conn)
			case strings.TrimSpace(message) == "554 5.5.1 Error: no valid recipients":
				sendResponse("QUIT\n", conn)
				fmt.Println(getCurrentTime() + " ERROR: No valid recipients found!")
			case strings.HasPrefix(message, "354"):
				sendResponse(emailData, conn)
				sendResponse("QUIT\n", conn)
				return true
			default:
				sendResponse("QUIT\n", conn)
				fmt.Println(getCurrentTime() + " ERROR: An unexpected response has occurred...")
				return false
			}
		}
	}
}
