package main

import (
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

func mailReceiver(socket net.Listener, configs map[string]string) {
	domainName := configs["EMAIL_DOMAIN"]
	debugMode := configs["DEBUG_MODE"]
	conn, err := socket.Accept()
	if err != nil {
		log.Fatal(err)
	}

	if _, err = conn.Write([]byte("220 LMTP Server Ready\n")); err != nil {
		log.Fatal(err)
	}
	fmt.Println(getCurrentTime() + "Initialising LMTP greeting")
	inData := false

	for {
		var size int
		buffer := make([]byte, 4096)

		size, err = conn.Read(buffer)
		if err != nil {
			fmt.Println(getCurrentTime() + "Error reading from LMTP greeting: " + err.Error())
			if err = conn.Close(); err != nil {
				log.Fatal(err)
			}
			return
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
				// TODO: Handle unknown email addresses
				// need db for this
				// for now, assume all email addresses are currently valid
				sendResponse("250 OK\n", conn)
			case strings.HasPrefix(message, "RCPT TO"):
				// TODO: Similar to MAIL FROM response, need to handle it correctly
				sendResponse("250 OK\n", conn)
			case strings.TrimSpace(message) == "DATA":
				sendResponse("354 Start mail input; end with <CRLF>.<CRLF>\n", conn)
				inData = true
			case inData:
				if strings.TrimSpace(message) == "." {
					inData = false
					mailSender(emailMessage, debugMode, configs)
					if _, err = conn.Write([]byte("250 OK (Sent to mailing list recipients)\n452 temporarily over quota\n")); err != nil {
						log.Fatal(err)
					}
				} else {
					emailMessage += message
				}
				fmt.Println(message)
			case strings.TrimSpace(message) == "QUIT":
				sendResponse(fmt.Sprintf("221 %s closing connection", domainName), conn)
				fmt.Println(getCurrentTime() + " S: Email successfully received, listening for more incoming emails")
				conn, err = socket.Accept()
				if _, err = conn.Write([]byte("220 LMTP Server Ready\n")); err != nil {
					log.Fatal(err)
				}
			default:
				if err = conn.Close(); err != nil {
					log.Fatal(err)
				}
				fmt.Println(getCurrentTime() + " ERROR: Unexpected response, closing connection...")
				conn, err = socket.Accept()
				if _, err = conn.Write([]byte("220 LMTP Server Ready\n")); err != nil {
					log.Fatal(err)
				}
			}
		}
	}

}

func mailSender(emailData string, debugMode string, configs map[string]string) {
	addr := configs["POSTFIX_ADDRESS"]
	port := configs["POSTFIX_PORT"]
	domainName := configs["EMAIL_DOMAIN"]

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

	isEnd := false
	for !isEnd {
		var size int
		buffer := make([]byte, 4096)

		size, err = conn.Read(buffer)
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
				// TODO: Remove example once db is set up and ready
				mailingList := "mailingList"
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
				isEnd = true
			default:
				sendResponse("QUIT\n", conn)
				fmt.Println(getCurrentTime() + " ERROR: An unexpected response has occurred...")
				isEnd = true
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

	socket, err := net.Listen("tcp", ":8024")
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(getCurrentTime() + " Starting up mail forwarding service...")
	mailReceiver(socket, configs)
}
