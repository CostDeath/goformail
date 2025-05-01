package mail

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"sync"
	"testing"
	"time"
)

var defaultEmail = model.Email{
	Id:          1,
	Rcpt:        []string{"list@domain.tld"},
	Sender:      "sender@domain.tld",
	Content:     "content",
	ReceivedAt:  time.Now(),
	NextRetry:   time.Now(),
	Exhausted:   3,
	Sent:        false,
	Approved:    true,
	ListMembers: []string{"user1@domain.tld", "user2@domain.tld"},
}

var defaultEmailLowerExhaust = model.Email{
	Id:          1,
	Rcpt:        []string{"list@domain.tld"},
	Sender:      "sender@domain.tld",
	Content:     "content",
	ReceivedAt:  time.Now(),
	NextRetry:   time.Now(),
	Exhausted:   2,
	Sent:        false,
	Approved:    true,
	ListMembers: []string{"user1@domain.tld", "user2@domain.tld"},
}

var defaultExhaustedEmail = model.Email{
	Id:          1,
	Rcpt:        []string{"list@domain.tld"},
	Sender:      "sender@domain.tld",
	Content:     "content",
	ReceivedAt:  time.Now(),
	NextRetry:   time.Now(),
	Exhausted:   0,
	Sent:        false,
	Approved:    true,
	ListMembers: []string{"user1@domain.tld", "user2@domain.tld"},
}

func TestCheckMail(t *testing.T) {
	mockDb := new(db.IDbMock)
	mockDb.On("GetAllReadyEmails").Return(&[]model.Email{defaultEmail})
	sender := NewEmailSender(nil, mockDb, util.MockConfigs)
	waitGroup := sync.WaitGroup{}

	sender.sendLock.Lock() // Prevent it from calling sender

	// Check it waits for queue lock
	go func() {
		sender.queueLock.Lock()
		time.Sleep(3 * time.Second)
		assert.Empty(t, sender.queue)
		sender.queueLock.Unlock()
	}()

	// Run function
	waitGroup.Add(1)
	go callFunctionWithWait(sender.CheckMail, &waitGroup)

	waitGroup.Wait()
	mockDb.AssertExpectations(t)
	assert.Equal(t, []model.Email{defaultEmail}, sender.queue)
}

func TestCheckMailIgnoredWhenLocked(t *testing.T) {
	mockDb := new(db.IDbMock)
	mockDb.On("GetAllReadyEmails").Panic("GetAllReadyEmails should not be called")
	sender := NewEmailSender(nil, mockDb, util.MockConfigs)
	waitGroup := sync.WaitGroup{}

	sender.checkLock.Lock() // Prevent it from calling the function

	// Run function
	waitGroup.Add(1)
	go callFunctionWithWait(sender.CheckMail, &waitGroup)

	waitGroup.Wait()
	assert.Empty(t, sender.queue)
}

func TestCheckMailDoesntAppendOnDbError(t *testing.T) {
	mockDb := db.NewIDbMockWithError(db.Unknown)
	mockDb.On("GetAllReadyEmails").Return(&[]model.Email{defaultEmail})
	sender := NewEmailSender(nil, mockDb, util.MockConfigs)
	waitGroup := sync.WaitGroup{}

	sender.checkLock.Lock() // Prevent it from calling the function

	// Run function
	waitGroup.Add(1)
	go callFunctionWithWait(sender.CheckMail, &waitGroup)

	waitGroup.Wait()
	assert.Empty(t, sender.queue)
}

func TestSendMail(t *testing.T) {
	mockMtp := new(IMtpHandlerMock)
	mockMtp.On("mailSender", defaultEmail.Sender, defaultEmail.ListMembers, defaultEmail.Content).Return(true)
	mockDb := new(db.IDbMock)
	mockDb.On("SetEmailAsSent", 1).Return()
	sender := NewEmailSender(mockMtp, mockDb, util.MockConfigs)
	waitGroup := sync.WaitGroup{}

	sender.checkLock.Lock() // Prevent it from calling checker

	// Check it waits for the first queue lock
	go func() {
		sender.queueLock.Lock()
		time.Sleep(3 * time.Second)
		sender.queue = []model.Email{defaultEmail, defaultEmail}
		sender.queueLock.Unlock()
	}()

	// Run function
	waitGroup.Add(1)
	time.Sleep(1 * time.Second)
	go callFunctionWithWait(sender.SendMail, &waitGroup)

	waitGroup.Wait()
	mockMtp.AssertExpectations(t)
	mockDb.AssertExpectations(t)
	mockMtp.AssertNumberOfCalls(t, "mailSender", 2)
	mockDb.AssertNumberOfCalls(t, "SetEmailAsSent", 2)
	assert.Empty(t, sender.queue)
}

func TestSendMailIgnoredWhenLocked(t *testing.T) {
	mockMtp := new(IMtpHandlerMock)
	mockMtp.On("mailSender", mock.Anything, mock.Anything, mock.Anything).Panic("GetAllReadyEmails should not be called")
	mockDb := new(db.IDbMock)
	mockDb.On("SetEmailAsSent", mock.Anything).Panic("GetAllReadyEmails should not be called")
	sender := NewEmailSender(mockMtp, mockDb, util.MockConfigs)
	waitGroup := sync.WaitGroup{}

	sender.sendLock.Lock() // Prevent it from calling the function

	// Run function
	waitGroup.Add(1)
	go callFunctionWithWait(sender.SendMail, &waitGroup)

	waitGroup.Wait()
	assert.Empty(t, sender.queue)
}

func TestSendMailSwapsOriginalSender(t *testing.T) {
	mockMtp := new(IMtpHandlerMock)
	mockMtp.On("mailSender", defaultEmail.Rcpt[0], defaultEmail.ListMembers, defaultEmail.Content).Return(true)
	mockDb := new(db.IDbMock)
	mockDb.On("SetEmailAsSent", 1).Return()
	sender := NewEmailSender(mockMtp, mockDb, map[string]string{"ORIGINAL_SENDER": "false", "CHECK_FREQUENCY": "60"})
	waitGroup := sync.WaitGroup{}

	sender.checkLock.Lock() // Prevent it from calling checker
	sender.queue = []model.Email{defaultEmail}

	// Run function
	waitGroup.Add(1)
	go callFunctionWithWait(sender.SendMail, &waitGroup)

	waitGroup.Wait()
	mockMtp.AssertExpectations(t)
	mockDb.AssertExpectations(t)
	assert.Empty(t, sender.queue)
}

func TestSendMailSetsRetry(t *testing.T) {
	mockMtp := new(IMtpHandlerMock)
	mockMtp.On("mailSender", defaultEmail.Sender, defaultEmail.ListMembers, defaultEmail.Content).Return(false)
	mockDb := new(db.IDbMock)
	mockDb.On("SetEmailRetry", mock.Anything).Run(func(args mock.Arguments) {
		actual := args.Get(0).(*model.Email)
		assert.Equal(t, 2, actual.Exhausted)
		assert.True(t, time.Now().Before(actual.NextRetry))
		assert.True(t, time.Now().Add(61*time.Minute).After(actual.NextRetry))
	}).Return()
	sender := NewEmailSender(mockMtp, mockDb, util.MockConfigs)
	waitGroup := sync.WaitGroup{}

	sender.checkLock.Lock() // Prevent it from calling checker
	sender.queue = []model.Email{defaultEmail}

	// Run function
	waitGroup.Add(1)
	go callFunctionWithWait(sender.SendMail, &waitGroup)

	waitGroup.Wait()
	mockMtp.AssertExpectations(t)
	mockDb.AssertExpectations(t)
	assert.Empty(t, sender.queue)
}

func TestSendMailSetsRetryDependentOnEmailAndConfig(t *testing.T) {
	mockMtp := new(IMtpHandlerMock)
	mockMtp.On("mailSender", defaultEmail.Sender, defaultEmail.ListMembers, defaultEmail.Content).Return(false)
	mockDb := new(db.IDbMock)
	mockDb.On("SetEmailRetry", mock.Anything).Run(func(args mock.Arguments) {
		actual := args.Get(0).(*model.Email)
		assert.Equal(t, 1, actual.Exhausted)
		assert.True(t, time.Now().Before(actual.NextRetry))
		assert.True(t, time.Now().Add(21*time.Minute).After(actual.NextRetry))
	}).Return()
	sender := NewEmailSender(mockMtp, mockDb, map[string]string{"ORIGINAL_SENDER": "true", "CHECK_FREQUENCY": "10"})
	waitGroup := sync.WaitGroup{}

	sender.checkLock.Lock() // Prevent it from calling checker
	sender.queue = []model.Email{defaultEmailLowerExhaust}

	// Run function
	waitGroup.Add(1)
	go callFunctionWithWait(sender.SendMail, &waitGroup)

	waitGroup.Wait()
	mockMtp.AssertExpectations(t)
	mockDb.AssertExpectations(t)
	assert.Empty(t, sender.queue)
}

func TestSendMailSetsRetryWhenExhausted(t *testing.T) {
	mockMtp := new(IMtpHandlerMock)
	mockMtp.On("mailSender", defaultEmail.Sender, defaultEmail.ListMembers, defaultEmail.Content).Return(false)
	mockDb := new(db.IDbMock)
	mockDb.On("SetEmailRetry", mock.Anything).Run(func(args mock.Arguments) {
		actual := args.Get(0).(*model.Email)
		assert.Equal(t, 0, actual.Exhausted)
		assert.Zero(t, actual.NextRetry)
	}).Return()
	sender := NewEmailSender(mockMtp, mockDb, util.MockConfigs)
	waitGroup := sync.WaitGroup{}

	sender.checkLock.Lock() // Prevent it from calling checker
	sender.queue = []model.Email{defaultExhaustedEmail}

	// Run function
	waitGroup.Add(1)
	go callFunctionWithWait(sender.SendMail, &waitGroup)

	waitGroup.Wait()
	mockMtp.AssertExpectations(t)
	mockDb.AssertExpectations(t)
	assert.Empty(t, sender.queue)
}

func TestCheckMailAndSendMailCallEachOther(t *testing.T) {
	mockMtp := new(IMtpHandlerMock)
	mockMtp.On("mailSender", defaultEmail.Sender, defaultEmail.ListMembers, defaultEmail.Content).Return(true)
	mockDb := new(db.IDbMock)
	mockDb.On("SetEmailAsSent", 1).Return()
	mockDb.On("GetAllReadyEmails").Return(&[]model.Email{defaultEmail}).Once()
	mockDb.On("GetAllReadyEmails").Return(&[]model.Email{}).Maybe()

	sender := NewEmailSender(mockMtp, mockDb, util.MockConfigs)
	waitGroup := sync.WaitGroup{}

	// Run function
	waitGroup.Add(1)
	go callFunctionWithWait(sender.CheckMail, &waitGroup)

	waitGroup.Wait()
	mockMtp.AssertExpectations(t)
	mockDb.AssertExpectations(t)
	mockMtp.AssertNumberOfCalls(t, "mailSender", 1)
	mockDb.AssertNumberOfCalls(t, "SetEmailAsSent", 1)
	mockDb.AssertNumberOfCalls(t, "GetAllReadyEmails", 2)
	assert.Empty(t, sender.queue)
}
