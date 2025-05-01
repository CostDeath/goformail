package mail

import (
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"log"
	"net"
	"strings"
	"time"
)

type EmailReceiver struct {
	mtp      *MtpHandler
	sender   *EmailSender
	socket   net.Listener
	lmtpPort string
}

func NewEmailReceiver(mtp *MtpHandler, sender *EmailSender, configs map[string]string) *EmailReceiver {
	return &EmailReceiver{
		mtp:      mtp,
		sender:   sender,
		lmtpPort: configs["LMTP_PORT"],
	}
}

func (receiver *EmailReceiver) Loop() {
	fmt.Println(getCurrentTime() + " Starting LMTP Service...")
	tcpSocket, err := createLMTPSocket(receiver.lmtpPort)
	if err != nil {
		log.Fatal(err)
	}
	receiver.socket = tcpSocket

	for {
		receiver.receiveMail()
	}
}

func (receiver *EmailReceiver) receiveMail() {
	conn, err := receiver.socket.Accept()
	if err != nil {
		log.Fatal(err)
	}

	// MAIL RECEIVER LOGIC
	email, err := receiver.mtp.mailReceiver(conn)
	if err != nil {
		fmt.Printf("%s %s\n", getCurrentTime(), err)
		if err = conn.Close(); err != nil {
			log.Fatal(err)
		}
		return
	}

	// QUEUE MAIL LOGIC
	mailForwardSuccess := false
	if email.Content != "" {
		for _, list := range email.Rcpt {
			id, approved, err := receiver.sender.db.GetApprovalFromListName(email.Sender, strings.Split(list, "@")[0])
			if err != nil {
				log.Println("Error getting list information: " + err.Err.Error())
			} else {
				emailPerList := model.Email{
					Rcpt:       []string{list},
					Sender:     email.Sender,
					Content:    email.Content,
					ReceivedAt: email.ReceivedAt,
					NextRetry:  time.Now(),
					Exhausted:  3,
					Sent:       false,
					Approved:   approved,
					List:       id,
				}
				err := receiver.sender.db.AddEmail(&emailPerList)
				if err != nil {
					log.Println(err.Err.Error())
				} else {
					mailForwardSuccess = true
				}
			}
		}

		if mailForwardSuccess {
			go receiver.sender.SendMail()
		}
	}

	// GOODBYE ACKNOWLEDGEMENT TO RESTART
	receiver.mtp.sendGoodbye(conn, mailForwardSuccess, email.RemainingAcks)
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
