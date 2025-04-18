package forwarding

import (
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func (e *emailCollectionError) Error() string {
	return fmt.Sprintf("%s: %s", e.errorType, e.Err.Error())
}

func LMTPService(configs map[string]string, db *db.Db) {
	fmt.Println(getCurrentTime() + " Starting LMTP Service...")
	lmtpPort, exists := configs["LMTP_PORT"]
	if !exists {
		log.Fatal("Missing LMTP_PORT config")
	}

	originalSender, exists := configs["ORIGINAL_SENDER"]
	if !exists {
		originalSender = "true"
		fmt.Println(getCurrentTime() + " S: ORIGINAL_SENDER config not found, setting it to true")
	}

	tcpSocket, err := createLMTPSocket(lmtpPort)
	if err != nil {
		log.Fatal(err)
	}

	bufferSizeConfig, exists := configs["BUFFER_SIZE"]
	if !exists {
		log.Fatal("Missing BUFFER_SIZE config")
	}
	bufferSize, err := strconv.Atoi(bufferSizeConfig)
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
				rcpt, err := db.GetRecipientsFromListName(strings.Split(mailingList, "@")[0])
				if err != nil {
					log.Println(err)
				}
				if originalSender == "true" {
					mailForwardSuccess = mailSender(data.from, rcpt, data, bufferSize, configs)
				} else {
					mailForwardSuccess = mailSender(mailingList, rcpt, data, bufferSize, configs)
				}
			}
		}

		// GOODBYE ACKNOWLEDGEMENT TO RESTART
		sendGoodbye(conn, mailForwardSuccess, data.remainingAcks)
	}
}

func getCurrentTime() string {
	t := time.Now()
	return fmt.Sprintf("[%d-%02d-%02d %02d:%02d:%02d]", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func createLMTPSocket(lmtpPort string) (net.Listener, error) {
	tcpSocket, err := net.Listen("tcp", fmt.Sprintf(":%s", lmtpPort))
	if err != nil {
		return nil, err
	}
	return tcpSocket, nil
}

func connectToSMTPSocket(smtpAddress string, smtpPort string) (net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", smtpAddress, smtpPort))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
