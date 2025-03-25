package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
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

func mailReceiver(socket net.Listener) {
	domainName := os.Getenv("EMAIL_DOMAIN")
	debugMode := os.Getenv("DEBUG_MODE")
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
					mailSender(emailMessage, debugMode)
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
			}
		}
	}

}

func mailSender(emailData string, debugMode string) {
	addr := os.Getenv("POSTFIX_ADDRESS")
	port := os.Getenv("POSTFIX_PORT")
	domainName := os.Getenv("EMAIL_DOMAIN")

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
				initial = false
			case strings.TrimSpace(message) == "250 CHUNKING":
				mailingList := "mailingList"
				sendResponse(fmt.Sprintf("MAIL FROM: %s@%s\n", mailingList, domainName), conn)
			case strings.HasPrefix(message, "250 2.1.0"):
				// TODO: Handle negative responses (?)
				recipients := []string{"sdk194", "dags"}
				for _, recipient := range recipients {
					sendResponse(fmt.Sprintf("RCPT TO: %s@%s\n", recipient, domainName), conn)
				}
				sendResponse("DATA\n", conn)
			case strings.TrimSpace(message) == "554 5.5.1 Error: no valid recipients":
				sendResponse("QUIT\n", conn)
				fmt.Println(getCurrentTime() + " ERROR: No valid recipients found!")
			case strings.HasPrefix(message, "354"):
				// TODO: Handle negative response
				sendResponse(emailData+"\n.\n", conn)
				sendResponse("QUIT\n", conn)
				isEnd = true
			}
		}
	}

}

func main() {
	/*
		socketDirectory := "/var/spool/postfix/goformail/"
		socketPath := fmt.Sprintf("%sgoformail_lmtp", socketDirectory)
		err := godotenv.Load("configs.env")
		if err != nil {
			log.Fatal(getCurrentTime() + "Error loading .env file")
		}

		err = os.MkdirAll(socketDirectory, 0755)

		err = os.Remove(socketPath) // remove existing socket when app restarts and socket exists
		if err == nil {
			fmt.Println(getCurrentTime() + " App has been restarted, recreating socket...")
		}

		socket, err := net.Listen("unix", socketPath)
		if err != nil {
			log.Fatal(err)
			return
		}


		defer func(socket net.Listener) {
			err = socket.Close()
			if err != nil {
				log.Fatal(getCurrentTime() + err.Error())
			}
		}(socket)


		if err = os.Chmod(socketPath, 0666); err != nil {
			log.Fatal(err)
			return
		}

	*/
	err := godotenv.Load("configs.env")
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
	mailReceiver(socket)
}
