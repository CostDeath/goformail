package mail

import (
	"github.com/stretchr/testify/mock"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"net"
)

type IMtpHandlerMock struct {
	mock.Mock
	IMtpHandler
}

func (mock *IMtpHandlerMock) sendGoodbye(conn net.Conn, mailForwardSuccess bool, remainingAcks []string) {
	mock.Called(conn, mailForwardSuccess, remainingAcks)
}

func (mock *IMtpHandlerMock) mailReceiver(conn net.Conn) (model.Email, error) {
	args := mock.Called(conn)
	if args.Get(1) == nil {
		return args.Get(0).(model.Email), nil
	}
	return args.Get(0).(model.Email), args.Get(1).(error)
}

func (mock *IMtpHandlerMock) mailSender(sender string, rcpt []string, content string) bool {
	args := mock.Called(sender, rcpt, content)
	return args.Bool(0)
}
