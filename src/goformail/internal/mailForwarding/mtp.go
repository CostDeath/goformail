package forwarding

import (
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

func validEmail(email string) bool {
	matches, err := regexp.Match(`^([a-z0-9\+\._\/&!][-a-z0-9\+\._\/&!]*)@(([a-z0-9][-a-z0-9]*\.)([-a-z0-9]+\.)*[a-z]{2,})$`, []byte(email))
	if err != nil {
		log.Fatal(err)
	}
	return matches
}

func mailReceiver(conn net.Conn, bufferSize int, configs map[string]string) (EmailData, error) {
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

				matches := validEmail(email)

				if !matches {
					sendResponse(fmt.Sprintf("550 5.1.1 <%s> user unknown", email), conn)
				} else {
					data.from = email
					sendResponse("250 OK\n", conn)
				}
			case strings.HasPrefix(message, "RCPT TO"):
				email := strings.Fields(message)[1][4:]
				email = email[:len(email)-1]
				if matches := validEmail(email); !matches {
					sendResponse(fmt.Sprintf("550 5.1.1 <%s> user unknown", email), conn)
				} else {
					data.rcpt = append(data.rcpt, email)
					sendResponse("250 OK\n", conn)
				}
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

func mailSender(mailingList string, emailData string, bufferSize int, configs map[string]string) bool {
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
				// TODO: Grab emails from db
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
