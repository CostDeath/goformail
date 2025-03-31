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
	if matches := validEmail(positive); !matches {
		t.Errorf("Email is not valid when it should be: %s", positive)
	}
}

func TestInvalidEmail(t *testing.T) {
	negative := "-this-doesnt@example.com"
	if matches := validEmail(negative); matches {
		t.Errorf("Email is valid when it should not be: %s", negative)
	}
}

func TestCreateLMTPSocket(t *testing.T) {
	tcpSocket, err := createLMTPSocket("8024")
	if err != nil {
		t.Error("tcpSocket was not created")
		return
	}

	err = tcpSocket.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func TestFailCreateLMTPSocket(t *testing.T) {
	tcpSocket, err := createLMTPSocket("not a port")
	if err == nil {
		t.Error("TCP socket was able to be created when it shouldn't have")
		if err = tcpSocket.Close(); err != nil {
			t.Error("Created tcp socket was not able to be closed")
		}
	}
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
	conn, err := connectToSMTPSocket("127.0.0.1", "8025")
	if err != nil {
		t.Error("Connection could not be established")
		conn, err = net.Dial("tcp", "127.0.0.1:8025")
		if err != nil {
			log.Fatal(err)
		}
		if err = conn.Close(); err != nil {
			log.Fatal(err)
		}
		waitGroup.Wait()
		return
	}

	err = conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	waitGroup.Wait()
}

func TestFailConnectToSMTP(t *testing.T) {
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
			t.Log("Goroutine within TestFailConnectToSMTP function finished")
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
	conn, err := connectToSMTPSocket("127.0.0.1", "81111")
	if err == nil {
		t.Error("Connection has been established when it shouldn't have been")
		waitGroup.Wait()
		return
	}

	// let server accept a client so it can move on
	conn, err = net.Dial("tcp", "127.0.0.1:8025")
	if err != nil {
		log.Fatal(err)
	}
	if err = conn.Close(); err != nil {
		log.Fatal(err)
	}
	waitGroup.Wait()
	return
}
