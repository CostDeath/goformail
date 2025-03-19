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
		var resp string
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

		for message := range messages {
			if debugMode == "true" {
				fmt.Print("POSTFIX: " + message)
			}
			switch {
			case strings.HasPrefix(message, "LHLO"):
				resp = fmt.Sprintf("250-%s\n250-PIPELINING\n250 SIZE\n", domainName)
				if _, err = conn.Write([]byte(resp)); err != nil {
					log.Fatal(err)
				}
				fmt.Println(getCurrentTime() + " S: " + resp)
			case strings.HasPrefix(message, "MAIL FROM"):
				// TODO: Handle unknown email addresses
				// for now, assume all email addresses are currently valid
				resp = fmt.Sprintf("250 OK\n")
				if _, err = conn.Write([]byte(resp)); err != nil {
					log.Fatal(err)
				}
				fmt.Println(getCurrentTime() + " S: " + resp)
			case strings.HasPrefix(message, "RCPT TO"):
				// TODO: Similar to MAIL FROM response, need to handle it correctly
				resp = fmt.Sprintf("250 OK\n")
				if _, err = conn.Write([]byte(resp)); err != nil {
					log.Fatal(err)
				}
				fmt.Println(getCurrentTime() + " S: " + resp)
			case strings.HasPrefix(message, "DATA"):
				resp = fmt.Sprintf("354 Start mail input; end with <CRLF>.<CRLF>\n")
				if _, err = conn.Write([]byte(resp)); err != nil {
					log.Fatal(err)
				}
				fmt.Println(getCurrentTime() + " S: " + resp)
				inData = true
			case inData:
				if strings.TrimSpace(message) == "." {
					inData = false
					if _, err = conn.Write([]byte("250 OK (Sent to mailing list recipients)\n452 temporarily over quota\n")); err != nil {
						log.Fatal(err)
					}
				}
				fmt.Println(message)
			case strings.TrimSpace(message) == "QUIT":
				if _, err = conn.Write([]byte(fmt.Sprintf("221 %s closing connection", domainName))); err != nil {
					log.Fatal(err)
				}
				fmt.Println(getCurrentTime() + " S: Email successfully received, listening for more incoming emails")
				conn, err = socket.Accept()
				if _, err = conn.Write([]byte("220 LMTP Server Ready\n")); err != nil {
					log.Fatal(err)
				}
			}
		}
	}

}

func main() {
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

	if err = os.Chmod(socketPath, 0666); err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(getCurrentTime() + " Starting up mail forwarding service...")
	mailReceiver(socket)
}
