package mail

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"log"
	"net"
	"regexp"
	"sync"
	"testing"
	"time"
)

var defaultEmailReceived = model.Email{
	Rcpt:          []string{"list1@example.domain", "list2@example.domain", "user@domain.tld"},
	Sender:        "sender@domain.tld",
	Content:       "content",
	ReceivedAt:    time.Now(),
	RemainingAcks: []string{"QUIT"},
}

func TestReceiveMail(t *testing.T) {
	mockMtp := new(IMtpHandlerMock)
	mockMtp.On("mailReceiver", mock.Anything).Return(defaultEmailReceived, nil)
	mockMtp.On("sendGoodbye", mock.Anything, true, defaultEmailReceived.RemainingAcks).Return(nil)
	mockSender := new(IEmailSenderMock)
	mockSender.On("SendMail").Return(nil)
	mockDb := new(db.IDbMock)
	mockDb.On("GetApprovalFromListName", "sender@domain.tld", "list1").Return(1, true)
	mockDb.On("GetApprovalFromListName", "sender@domain.tld", "list2").Return(2, false)
	mockDb.On("AddEmail", mock.MatchedBy(func(e *model.Email) bool {
		expected := createEmail("list1@example.domain", 1, true)
		return e.Rcpt[0] == expected.Rcpt[0] && e.Sender == expected.Sender && e.Content == expected.Content &&
			e.ReceivedAt.Equal(expected.ReceivedAt) && e.Exhausted == expected.Exhausted && e.Sent == expected.Sent &&
			e.List == expected.List && e.ReceivedAt.Before(e.NextRetry)
	})).Return()
	mockDb.On("AddEmail", mock.MatchedBy(func(e *model.Email) bool {
		expected := createEmail("list2@example.domain", 2, false)
		return e.Rcpt[0] == expected.Rcpt[0] && e.Sender == expected.Sender && e.Content == expected.Content &&
			e.ReceivedAt.Equal(expected.ReceivedAt) && e.Exhausted == expected.Exhausted && e.Sent == expected.Sent &&
			e.List == expected.List && e.ReceivedAt.Before(e.NextRetry)
	})).Return()

	receiver := NewEmailReceiver(mockMtp, mockSender, mockDb, util.MockConfigs)
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer func(listener net.Listener) {
		err := listener.Close()
		require.NoError(t, err)
	}(listener)
	receiver.socket = listener
	waitGroup := sync.WaitGroup{}

	// Run function
	waitGroup.Add(1)
	go callFunctionWithWait(receiver.receiveMail, &waitGroup)
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", listener.Addr().(*net.TCPAddr).Port))
	require.NoError(t, err)
	defer func(conn net.Conn) {
		err := conn.Close()
		require.NoError(t, err)
	}(conn)

	waitGroup.Wait()
	//mockMtp.AssertExpectations(t)
	//mockSender.AssertExpectations(t)
	//mockDb.AssertExpectations(t)
	mockDb.AssertNotCalled(t, "GetApprovalFromListName", "sender@domain.tld", "user")
}

func TestReceiveMailEndsWhenMailReceiverErrors(t *testing.T) {
	mockMtp := new(IMtpHandlerMock)
	mockMtp.On("mailReceiver", mock.Anything).Return(defaultEmailReceived, &emailCollectionError{"READ_ERROR", nil})
	mockMtp.On("sendGoodbye", mock.Anything, false, defaultEmailReceived.RemainingAcks).Return(nil)

	receiver := NewEmailReceiver(mockMtp, nil, nil, util.MockConfigs)
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer func(listener net.Listener) {
		err := listener.Close()
		require.NoError(t, err)
	}(listener)
	receiver.socket = listener
	waitGroup := sync.WaitGroup{}

	// Run function
	waitGroup.Add(1)
	go callFunctionWithWait(receiver.receiveMail, &waitGroup)
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", listener.Addr().(*net.TCPAddr).Port))
	require.NoError(t, err)
	defer func(conn net.Conn) {
		err := conn.Close()
		require.NoError(t, err)
	}(conn)

	waitGroup.Wait()
	mockMtp.AssertExpectations(t)
}

func TestReceiveMailDoesNotSendWhenAllDbError(t *testing.T) {
	mockMtp := new(IMtpHandlerMock)
	mockMtp.On("mailReceiver", mock.Anything).Return(defaultEmailReceived, nil)
	mockMtp.On("sendGoodbye", mock.Anything, false, defaultEmailReceived.RemainingAcks).Return(nil)
	mockDb := db.NewIDbMockWithError(db.Unknown)
	mockDb.On("GetApprovalFromListName", "sender@domain.tld", "list1").Return(0, false, nil)
	mockDb.On("GetApprovalFromListName", "sender@domain.tld", "list2").Return(0, false, nil)

	receiver := NewEmailReceiver(mockMtp, nil, mockDb, util.MockConfigs)
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer func(listener net.Listener) {
		err := listener.Close()
		require.NoError(t, err)
	}(listener)
	receiver.socket = listener
	waitGroup := sync.WaitGroup{}

	// Run function
	waitGroup.Add(1)
	go callFunctionWithWait(receiver.receiveMail, &waitGroup)
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", listener.Addr().(*net.TCPAddr).Port))
	require.NoError(t, err)
	defer func(conn net.Conn) {
		err := conn.Close()
		require.NoError(t, err)
	}(conn)

	waitGroup.Wait()
	mockMtp.AssertExpectations(t)
	mockDb.AssertExpectations(t)
	mockDb.AssertNotCalled(t, "AddEmail", mock.Anything)
}

func TestGetCurrentTime(t *testing.T) {
	formattedTime := getCurrentTime()

	// EXPECTED FORMAT [YYYY-MM-DD HH:MM:SS]

	matched, err := regexp.Match(
		`\[([0-9]{4})(-([0-9]{2})){2}\s([0-5][0-9]:){2}[0-5][0-9]]`,
		[]byte(formattedTime),
	)
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
		ConnectSMTPSocketMock(waitGroup)
		t.Log("Goroutine function has finished")
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
		ConnectFailSMTPSocketMock(waitGroup)
		t.Log("Goroutine within TestFailConnectToSMTP function finished")
	}()

	waitGroup.Wait()
	waitGroup.Add(1)
	conn, err := connectToSMTPSocket("127.0.0.1", "81111")
	if err == nil {
		t.Error("Connection has been established when it shouldn't have been")
		waitGroup.Wait()
		if err = conn.Close(); err != nil {
			t.Error("There was an error when attempting to close the connection " + err.Error())
		}
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
}

func createEmail(recipient string, list int, approved bool) *model.Email {
	return &model.Email{
		Rcpt:          []string{recipient},
		Sender:        defaultEmailReceived.Sender,
		Content:       defaultEmailReceived.Content,
		ReceivedAt:    defaultEmailReceived.ReceivedAt,
		NextRetry:     time.Time{},
		Exhausted:     3,
		Sent:          false,
		Approved:      approved,
		RemainingAcks: defaultEmailReceived.RemainingAcks,
		List:          list,
	}
}
