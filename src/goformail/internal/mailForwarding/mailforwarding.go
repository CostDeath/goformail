package mailfowarding

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net"
	"strings"
	"time"
)

func getCurrentTime() string {
	t := time.Now()
	return fmt.Sprintf("[%d-%02d-%02d %02d:%02d:%02d]", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
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

func LMTPService(configs map[string]string) {
	lmtpPort := configs["LMTP_PORT"]

	tcpSocket, err := net.Listen("tcp", fmt.Sprintf(":%s", lmtpPort))
	if err != nil {
		log.Fatal(err)
		return
	}

	var conn net.Conn
	var mailForwardSuccess bool
	for {
		conn, err = tcpSocket.Accept()
		if err != nil {
			log.Fatal(err)
		}
		mailForwardSuccess = false

		// MAIL RECEIVER LOGIC
		data := MailReceiver(conn, configs)
		if _, containsError := data["READ_ERROR"]; containsError {
			fmt.Println(getCurrentTime() + "Error reading from LMTP greeting: " + err.Error())
			if err = conn.Close(); err != nil {
				log.Fatal(err)
			}
			continue // want to go back to loop
		}
		if _, containsError := data["RESPONSE_ERROR"]; containsError {
			fmt.Println(getCurrentTime() + " ERROR: Unexpected response, closing connection...")
			if err = conn.Close(); err != nil {
				log.Fatal(err)
			}
			continue // want to go back to loop
		}
		// SEND MAIL LOGIC
		if emailData, exists := data["EMAIL_DATA"]; exists {
			mailingLists := strings.Fields(data["RCPTS"])

			for _, mailingList := range mailingLists {
				mailForwardSuccess = MailSender(mailingList, emailData, configs)
			}
		}
		// GOODBYE ACKNOWLEDGEMENT TO RESTART
		sendGoodbye(conn, mailForwardSuccess, configs["REMAINING_ACK"])
	}
}

func MailReceiver(conn net.Conn, configs map[string]string) map[string]string {
	domainName := configs["EMAIL_DOMAIN"]
	debugMode := configs["DEBUG_MODE"]

	result := make(map[string]string)
	result["RCPTS"] = ""
	result["REMAINING_ACK"] = ""

	if _, err := conn.Write([]byte("220 LMTP Server Ready\n")); err != nil {
		log.Fatal(err)
	}
	fmt.Println(getCurrentTime() + "Initialising LMTP greeting")
	inData := false
	for {
		var size int
		buffer := make([]byte, 4096)

		size, err := conn.Read(buffer)
		if err != nil {
			result["READ_ERROR"] = err.Error()
			return result
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
				// TODO: Handle permissions to be able to send from mailing lists
				// need db for this
				// for now, assume all email addresses are currently valid
				sendResponse("250 OK\n", conn)
			case strings.HasPrefix(message, "RCPT TO"):
				// TODO: Handle possible cases where local recipient table and db are not fully up to date with each other
				// E.G. Mailing list email was deleted in db but local recipient table still has it mapped to the LMTP port
				email := strings.Fields(message)[1][4:]
				email = email[:len(email)-1]
				result["RCPTS"] += fmt.Sprintf(" %s", email)
				sendResponse("250 OK\n", conn)
			case strings.TrimSpace(message) == "DATA":
				sendResponse("354 Start mail input; end with <CRLF>.<CRLF>\n", conn)
				inData = true
			case inData:
				if strings.TrimSpace(message) == "." {
					inData = false
					result["EMAIL_DATA"] = emailMessage
				} else {
					emailMessage += message
				}
				fmt.Println(message)
			case strings.TrimSpace(message) == "QUIT":
				result["REMAINING_ACK"] += message
				return result
			default:
				result["RESPONSE_ERROR"] = message
				return result
			}
		}
	}

}

func MailSender(mailingList string, emailData string, configs map[string]string) bool {
	addr := configs["POSTFIX_ADDRESS"]
	port := configs["POSTFIX_PORT"]
	domainName := configs["EMAIL_DOMAIN"]
	debugMode := configs["DEBUG_MODE"]
	timeoutDuration, err := time.ParseDuration(configs["TIMEOUT_DURATION"] + "s")
	if err != nil {
		fmt.Println(getCurrentTime() + " ERROR: Could not parse timeout duration: " + err.Error())
		return false
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", addr, port))
	if err != nil {
		log.Fatal(getCurrentTime() + err.Error())
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
		buffer := make([]byte, 4096)

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
				sendResponse(emailData+"\n.\n", conn)
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

func main() {
	configs, err := godotenv.Read("configs.cf")
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(getCurrentTime() + " Starting up LMTP service...")
	LMTPService(configs)
}
