package model

import "time"

type Email struct {
	Id            int       `json:"id"`
	Rcpt          []string  `json:"rcpt"`
	Sender        string    `json:"sender"`
	Content       string    `json:"content"`
	ReceivedAt    time.Time `json:"received_at"`
	NextRetry     time.Time `json:"next_retry"`
	Exhausted     int       `json:"exhausted"`
	Sent          bool      `json:"sent"`
	ListName      string    `json:"list"`
	Approved      bool      `json:"approved"`
	RemainingAcks []string
	List          int
	ListMembers   []string
}

type EmailReqs struct {
	Offset          int  `json:"offset"`
	List            *int `json:"list"`
	Archived        bool `json:"archived"`
	Exhausted       bool `json:"exhausted"`
	PendingApproval bool `json:"pending_approval"`
}

type EmailResponse struct {
	Offset int     `json:"offset"`
	Emails []Email `json:"emails"`
}
