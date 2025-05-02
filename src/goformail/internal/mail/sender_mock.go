package mail

import (
	"github.com/stretchr/testify/mock"
	"log"
)

type IEmailSenderMock struct {
	mock.Mock
	IEmailSender
}

func (mock *IEmailSenderMock) CheckMail() {
	mock.Called()
}

func (mock *IEmailSenderMock) SendMail() {
	args := mock.Called()
	log.Print(args) // force mock to acknowledge the call
}
