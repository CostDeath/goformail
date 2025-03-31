package forwarding

import (
	"log"
	"net"
	"regexp"
	"sync"
	"testing"
)

func TestGetCurrentTime(t *testing.T) {
	formattedTime := getCurrentTime()

	// EXPECTED FORMAT [YYYY-MM-DD HH:MM:SS]
	matched, err := regexp.Match(`\[([0-9]{4})(-([0-9]{2})){2}\s([0-5][0-9]:){2}[0-5][0-9]]`, []byte(formattedTime))
	if err != nil {
		log.Fatal(err)
	}
	if !matched {
		t.Errorf("Not the expected format of [YYYY-MM-DD HH:mm:ss], got %s", formattedTime)
	}
}

func TestValidEmail(t *testing.T) {
	positive := "this-works@example.com"
	negative := "-this-doesnt@example.com"
	if matches := validEmail(positive); !matches {
		t.Errorf("Email is not valid when it should be: %s", positive)
	}
	if matches := validEmail(negative); matches {
		t.Errorf("Email is valid when it should not be: %s", negative)
	}
}

func TestConnectToLMTP(t *testing.T) {
	tcpSocket := createLMTPSocket("8024")
	if tcpSocket == nil {
		t.Errorf("tcpSocket is nil")
		return
	}

	defer func(tcpSocket net.Listener) {
		err := tcpSocket.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(tcpSocket)
}

func TestConnectToSMTP(t *testing.T) {
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)
	// MOCK Listener
	go func() {
		tcpSocket, err := net.Listen("tcp", "127.0.0.1:8025")
		if err != nil {
			log.Fatal(err)
		}
		defer func(tcpSocket net.Listener) {
			err = tcpSocket.Close()
			if err != nil {
				log.Fatal(err)
			}
			t.Log("Goroutine within testConnectToSMTP function finished")
			waitGroup.Done()
		}(tcpSocket)

		waitGroup.Done()
		_, err = tcpSocket.Accept()
		if err != nil {
			log.Fatal(err)
		}
	}()

	waitGroup.Wait()
	waitGroup.Add(1)
	conn := connectToSMTPSocket("127.0.0.1", "8025")

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
		waitGroup.Wait()
	}(conn)
}
