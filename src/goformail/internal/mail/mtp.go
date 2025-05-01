package mail

import (
	"errors"
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type emailCollectionError struct {
	errorType string
	Err       error
}

func (e *emailCollectionError) Error() string {
	return fmt.Sprintf("%s: %s", e.errorType, e.Err.Error())
}

type IMtpHandler interface {
	sendGoodbye(conn net.Conn, mailForwardSuccess bool, remainingAcks []string)
	mailReceiver(conn net.Conn) (model.Email, error)
	mailSender(sender string, rcpt []string, content string) bool
}

type MtpHandler struct {
	IMtpHandler
	bufferSize      int
	domain          string
	postfixAddr     string
	postfixPort     string
	timeoutDuration time.Duration
	debugMode       string
}

func NewMtpHandler(configs map[string]string) *MtpHandler {
	bufferSize, _ := strconv.Atoi(configs["BUFFER_SIZE"])
	timeoutDuration, _ := time.ParseDuration(configs["TIMEOUT_DURATION"] + "s")
	return &MtpHandler{
		bufferSize:      bufferSize,
		domain:          configs["EMAIL_DOMAIN"],
		postfixAddr:     configs["POSTFIX_ADDRESS"],
		postfixPort:     configs["POSTFIX_PORT"],
		timeoutDuration: timeoutDuration,
		debugMode:       configs["DEBUG_MODE"],
	}
}

func (mtp *MtpHandler) sendGoodbye(conn net.Conn, mailForwardSuccess bool, remainingAcks []string) {
	if mailForwardSuccess {
		sendResponse("250 OK (Email was successfully queued)\n452 temporarily over quota\n", conn)
	} else {
		sendResponse("250 OK (However, email was not queued)\n452 temporarily over quota\n", conn)
	}

	for _, message := range remainingAcks {
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

func validEmail(email string) bool {
	matches, err := regexp.Match(
		`^([A-z0-9+._/&!][-A-z0-9+._/&!]*)@(([a-z0-9][-a-z0-9]*\.)([-a-z0-9]+\.)*[a-z]{2,})$`,
		[]byte(email),
	)
	if err != nil {
		log.Fatal(err)
	}
	return matches
}

func sendResponse(resp string, conn net.Conn) {
	if _, err := conn.Write([]byte(resp)); err != nil {
		log.Fatal(err)
	}
	fmt.Println(getCurrentTime() + " S: " + resp)
}

func (mtp *MtpHandler) mailReceiver(conn net.Conn) (model.Email, error) {
	data := model.Email{ReceivedAt: time.Now()}

	if _, err := conn.Write([]byte("220 LMTP Server Ready\n")); err != nil {
		log.Fatal(err)
	}
	fmt.Println(getCurrentTime() + "Initialising LMTP greeting")
	inData := false
	fullStopFound := false
	var emailMessage string
	for {
		var size int
		buffer := make([]byte, mtp.bufferSize)

		size, err := conn.Read(buffer)
		if err != nil {
			return data, &emailCollectionError{"READ_ERROR", err}
		}

		messages := strings.Lines(string(buffer[:size]))

		for message := range messages {
			if mtp.debugMode == "true" {
				fmt.Print("POSTFIX: " + message)
			}
			switch {
			case inData:
				if strings.TrimSpace(message) == "." {
					fullStopFound = true
					emailMessage += message
				} else if strings.TrimSpace(message) == "QUIT" && fullStopFound {
					data.Content = emailMessage
					data.RemainingAcks = append(data.RemainingAcks, message)
					return data, nil
				} else {
					emailMessage += message
					// need to reset to false if fullstop was found earlier but QUIT message was not followed up next
					fullStopFound = false
				}
			case strings.HasPrefix(message, "LHLO"):
				sendResponse(fmt.Sprintf("250-%s\n250-PIPELINING\n250 SIZE\n", mtp.domain), conn)
			case strings.HasPrefix(message, "MAIL FROM"):
				email := strings.Fields(message)[1][6:]
				email = email[:len(email)-1]

				matches := validEmail(email)

				if !matches {
					sendResponse(fmt.Sprintf("550 5.1.1 <%s> user unknown", email), conn)
				} else {
					data.Sender = email
					sendResponse("250 OK\n", conn)
				}
			case strings.HasPrefix(message, "RCPT TO"):
				email := strings.Fields(message)[1][4:]
				email = email[:len(email)-1]
				if matches := validEmail(email); !matches {
					sendResponse(fmt.Sprintf("550 5.1.1 <%s> user unknown", email), conn)
				} else {
					data.Rcpt = append(data.Rcpt, email)
					sendResponse("250 OK\n", conn)
				}
			case strings.TrimSpace(message) == "DATA":
				sendResponse("354 Start mail input; end with <CRLF>.<CRLF>\n", conn)
				inData = true

			default:
				sendResponse("500 Error: command not recognised\n", conn)
				return data, &emailCollectionError{"UNEXPECTED_RESPONSE_ERROR", errors.New(message)}
			}
		}
	}

}

func (mtp *MtpHandler) mailSender(sender string, rcpt []string, content string) bool {
	conn, err := connectToSMTPSocket(mtp.postfixAddr, mtp.postfixPort)
	if err != nil {
		log.Fatal(err)
	}

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
		buffer := make([]byte, mtp.bufferSize)

		err = conn.SetDeadline(time.Now().Add(mtp.timeoutDuration)) // Time out after 5 seconds

		size, err = conn.Read(buffer)

		// handle timeout
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			fmt.Println(getCurrentTime() + " S: MTA not responding, is the correct port configured?")
			return false
		}

		messages := strings.Lines(string(buffer[:size]))

		for message := range messages {
			if mtp.debugMode == "true" {
				fmt.Print("POSTFIX: " + message)
			}
			switch {
			case initial:
				sendResponse(fmt.Sprintf("EHLO %s\n", mtp.domain), conn)

				size, err = conn.Read(buffer)
				if err != nil {
					log.Fatal(err)
				}

				messages = strings.Lines(string(buffer[:size]))
				for message = range messages {
					if !strings.HasPrefix(message, "250") {
						fmt.Println(getCurrentTime() + " S: Unsupported commands, not sending email...")
						sendResponse("QUIT", conn)
						return false
					}
				}

				sendResponse(fmt.Sprintf("MAIL FROM: %s\n", sender), conn)

				initial = false
			case strings.HasPrefix(message, "250 2.1.0"):
				recipients := rcpt
				for _, recipient := range recipients {
					sendResponse(fmt.Sprintf("RCPT TO: %s\n", recipient), conn)
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
				sendResponse(content, conn)
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
