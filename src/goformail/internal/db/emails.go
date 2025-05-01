package db

import (
	"github.com/lib/pq"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
)

func (db *Db) GetAllReadyEmails() (*[]model.Email, *Error) {
	rows, err := db.conn.Query(`
		SELECT emails.id, emails.rcpt, emails.sender, emails.content, emails.exhausted, lists.recipients
		FROM emails JOIN lists ON emails.list = lists.id
		WHERE emails.sent = false AND emails.approved = true AND emails.next_retry < NOW() AND emails.exhausted > 0;
	`)
	if err != nil {
		return nil, getError(err)
	}

	var emails []model.Email
	for rows.Next() {
		email := model.Email{}
		if err := rows.Scan(&email.Id, pq.Array(&email.Rcpt), &email.Sender, &email.Content, &email.Exhausted,
			pq.Array(&email.ListMembers)); err != nil {
			return nil, getError(err)
		}

		emails = append(emails, email)
	}

	return &emails, nil
}

func (db *Db) AddEmail(email *model.Email) *Error {
	if _, err := db.conn.Exec(`
		INSERT INTO emails (rcpt, sender, content, received_at, next_retry, exhausted, sent, list, approved)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id
	`, pq.Array(email.Rcpt), email.Sender, email.Content, email.ReceivedAt, email.NextRetry,
		email.Exhausted, email.Sent, email.List, email.Approved); err != nil {
		return getError(err)
	}
	return nil
}

func (db *Db) SetEmailAsSent(id int) *Error {
	res, err := db.conn.Exec(`
		UPDATE emails SET sent = true WHERE id = $1;
	`, id)
	if err != nil {
		return getError(err)
	} else if count, err := res.RowsAffected(); count == 0 || err != nil {
		return &Error{Err: err, Code: ErrNoRows}
	}

	return nil
}

func (db *Db) SetEmailRetry(email *model.Email) *Error {
	res, err := db.conn.Exec(`
		UPDATE emails SET next_retry = $1, exhausted = $2 WHERE id = $3;
	`, email.NextRetry, email.Exhausted, email.Id)
	if err != nil {
		return getError(err)
	} else if count, err := res.RowsAffected(); count == 0 || err != nil {
		return &Error{Err: err, Code: ErrNoRows}
	}

	return nil
}
