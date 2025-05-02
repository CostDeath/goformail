package mail

import (
	"github.com/stretchr/testify/mock"
)

type IEmailSenderMock struct {
	mock.Mock
	IEmailSender
}

func (mock *IEmailSenderMock) CheckMail() {
	mock.Called()
}

func (mock *IEmailSenderMock) SendMail() {
	mock.Called()
}
