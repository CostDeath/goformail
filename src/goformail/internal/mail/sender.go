package mail

import (
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"log"
	"strconv"
	"sync"
	"time"
)

type IEmailSender interface {
	CheckMail()
	SendMail()
}

type EmailSender struct {
	IEmailSender
	mtp            IMtpHandler
	db             db.IDb
	queue          []model.Email
	checkLock      sync.Mutex
	sendLock       sync.Mutex
	queueLock      sync.Mutex
	originalSender bool
	checkFrequency time.Duration
}

func NewEmailSender(mtp IMtpHandler, db db.IDb, configs map[string]string) *EmailSender {
	originalSender, _ := strconv.ParseBool(configs["ORIGINAL_SENDER"])
	checkFrequency, _ := time.ParseDuration(configs["CHECK_FREQUENCY"] + "m")
	return &EmailSender{
		mtp:            mtp,
		db:             db,
		queue:          []model.Email{},
		checkLock:      sync.Mutex{},
		sendLock:       sync.Mutex{},
		queueLock:      sync.Mutex{},
		originalSender: originalSender,
		checkFrequency: checkFrequency,
	}
}

func (sender *EmailSender) Loop() {
	lastCheck := time.Now()
	for {
		time.Sleep(sender.checkFrequency - time.Since(lastCheck))
		lastCheck = time.Now()

		sender.SendMail()
	}
}

func (sender *EmailSender) CheckMail() {
	if sender.checkLock.TryLock() {
		emails, err := sender.db.GetAllReadyEmails()
		if err != nil {
			log.Println("Error fetching ready emails: ", err.Err.Error())
			sender.checkLock.Unlock()
			return
		}
		if len(*emails) > 0 {
			sender.queueLock.Lock()
			sender.queue = append(sender.queue, *emails...)
			sender.queueLock.Unlock()
			sender.checkLock.Unlock()
			sender.SendMail()
		} else {
			sender.checkLock.Unlock()
		}
	}
}

func (sender *EmailSender) SendMail() {
	if sender.sendLock.TryLock() {
		sender.queueLock.Lock()
		for len(sender.queue) > 0 {
			email := sender.queue[0]
			sender.queueLock.Unlock()

			if !sender.originalSender {
				email.Sender = email.Rcpt[0]
			}
			success := sender.mtp.mailSender(email.Sender, email.ListMembers, email.Content)

			if success {
				err := sender.db.SetEmailAsSent(email.Id)
				if err != nil {
					log.Println("Error setting email as sent: ", err, email)
				}
			} else if email.Exhausted > 0 {
				email.NextRetry = time.Now().Add(sender.checkFrequency * time.Duration(4-email.Exhausted))
				email.Exhausted--
				err := sender.db.SetEmailRetry(&email)
				if err != nil {
					log.Println("Error setting email retry: ", err, email)
				}
			} else {
				email.NextRetry = time.Time{}
				email.Exhausted = 0
				err := sender.db.SetEmailRetry(&email)
				if err != nil {
					log.Println("Error setting email retry: ", err, email)
				}
			}

			sender.queueLock.Lock()
			sender.queue = sender.queue[1:]
		}
		sender.queueLock.Unlock()
		sender.sendLock.Unlock()
		sender.CheckMail()
	}
}
