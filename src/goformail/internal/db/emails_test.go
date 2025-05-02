package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"testing"
	"time"
)

var loc, _ = time.LoadLocation("Etc/UTC")

type DbEmailsSuite struct {
	suite.Suite
	c  *postgres.PostgresContainer
	db *Db
}

func TestDbEmailsSuite(t *testing.T) {
	suite.Run(t, new(DbEmailsSuite))
}

func (suite *DbEmailsSuite) SetupSuite() {
	ctx := context.Background()
	c, err := postgres.Run(
		ctx,
		"postgres:latest",
		postgres.WithDatabase("goformail"),
		postgres.WithUsername("goformail"),
		postgres.WithPassword("password"),
		postgres.WithInitScripts("scripts/emails_test.sql"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp")),
		testcontainers.WithHostConfigModifier(func(hostConfig *container.HostConfig) {
			hostConfig.AutoRemove = true
			hostConfig.Tmpfs = map[string]string{"/var/lib/postgresql/data": "rw"}
		}),
	)
	suite.Require().NoError(err)

	host, _ := c.Host(ctx)
	mappedPort, _ := c.MappedPort(ctx, "5432")
	info := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, mappedPort.Port(), util.MockConfigs["SQL_USER"], util.MockConfigs["SQL_PASSWORD"],
		util.MockConfigs["SQL_DB_NAME"])
	db, err := sql.Open("postgres", info)
	suite.Require().NoError(err)

	suite.c, suite.db = c, &Db{conn: db, batchSize: 3}
}

func (suite *DbEmailsSuite) TearDownSuite() {
	if err := testcontainers.TerminateContainer(suite.c); err != nil {
		suite.T().Log(err)
	}
}

func (suite *DbEmailsSuite) TestGetAllReadyEmails() {
	// Run function
	actual, actualErr := suite.db.GetAllReadyEmails()

	// Get expected
	rows, err := suite.db.conn.Query(`
		SELECT emails.id, emails.rcpt, emails.sender, emails.content, emails.exhausted, lists.recipients
		FROM emails JOIN lists ON emails.list = lists.id
		WHERE emails.sent = false AND emails.approved = true AND emails.next_retry < NOW() AND emails.exhausted > 0;
	`)
	suite.Require().NoError(err)

	var expected []model.Email
	for rows.Next() {
		email := model.Email{}
		err := rows.Scan(&email.Id, pq.Array(&email.Rcpt), &email.Sender, &email.Content, &email.Exhausted,
			pq.Array(&email.ListMembers))
		suite.Require().NoError(err)

		expected = append(expected, email)
	}

	suite.Nil(actualErr)
	suite.Len(*actual, 1)
	suite.Equal(&expected, actual)
}

func (suite *DbEmailsSuite) TestAddEmail() {
	// Run function
	expected := createEmail(true, true, 3)
	actualErr := suite.db.AddEmail(expected)

	// Check list was created properly
	var actual model.Email
	err := suite.db.conn.QueryRow(`
		SELECT rcpt, sender, content, received_at, next_retry, exhausted, sent, list, approved FROM emails WHERE list = $1
	`, 3).Scan(pq.Array(&actual.Rcpt), &actual.Sender, &actual.Content, &actual.ReceivedAt, &actual.NextRetry,
		&actual.Exhausted, &actual.Sent, &actual.List, &actual.Approved)
	suite.Require().NoError(err)

	suite.Nil(actualErr)
	suite.Equal(expected, &actual)
	suite.True(time.Now().After(actual.ReceivedAt))
}

func (suite *DbEmailsSuite) TestSetEmailAsSent() {
	// Run function
	actualErr := suite.db.SetEmailAsSent(5)

	// Check email was updated properly
	var sent bool
	err := suite.db.conn.QueryRow(`
		SELECT sent FROM emails WHERE id = $1
	`, 5).Scan(&sent)
	suite.Require().NoError(err)

	suite.Nil(actualErr)
	suite.True(sent)
}

func (suite *DbEmailsSuite) TestSetEmailAsSentReturnsErrorWhenNoEmail() {
	// Run function
	actualErr := suite.db.SetEmailAsSent(0)

	suite.NotNil(actualErr)
	suite.Equal(ErrNoRows, actualErr.Code)
}

func (suite *DbEmailsSuite) TestSetEmailRetry() {
	// Run function
	expected := createEmail(false, true, 2)
	expected.Id = 5
	expected.NextRetry = expected.NextRetry.Round(time.Microsecond)
	actualErr := suite.db.SetEmailRetry(expected)

	// Check email was updated properly
	actual := model.Email{}
	err := suite.db.conn.QueryRow(`
		SELECT next_retry, exhausted FROM emails WHERE id = $1
	`, 5).Scan(&actual.NextRetry, &actual.Exhausted)
	suite.Require().NoError(err)

	suite.Nil(actualErr)
	suite.Equal(expected.NextRetry, actual.NextRetry)
	suite.Equal(expected.Exhausted, actual.Exhausted)
}

func (suite *DbEmailsSuite) TestSetEmailRetryErrorWhenNoEmail() {
	// Run function
	email := createEmail(false, true, 2)
	actualErr := suite.db.SetEmailRetry(email)

	suite.NotNil(actualErr)
	suite.Equal(ErrNoRows, actualErr.Code)
}

func (suite *DbEmailsSuite) TestSetEmailAsApproved() {
	// Run function
	actualErr := suite.db.SetEmailAsApproved(4)

	// Check email was updated properly
	var approved bool
	err := suite.db.conn.QueryRow(`
		SELECT approved FROM emails WHERE id = $1
	`, 4).Scan(&approved)
	suite.Require().NoError(err)

	suite.Nil(actualErr)
	suite.True(approved)
}

func (suite *DbEmailsSuite) TestSetEmailAsApprovedReturnsErrorWhenNoEmail() {
	// Run function
	actualErr := suite.db.SetEmailAsApproved(0)

	suite.NotNil(actualErr)
	suite.Equal(ErrNoRows, actualErr.Code)
}

func createEmail(sent bool, approved bool, list int) *model.Email {
	return &model.Email{
		Rcpt:       []string{"test-gen@test.tld"},
		Sender:     "sender@test-0.tld",
		Content:    "content",
		ReceivedAt: time.Now().In(loc).Round(time.Microsecond),
		NextRetry:  time.Now().In(loc).Round(time.Microsecond),
		Exhausted:  2,
		Sent:       approved,
		Approved:   sent,
		List:       list,
	}
}

func (suite *DbEmailsSuite) TestGetAllEmails() {
	// Run function
	actual, actualErr := suite.db.GetAllEmails(&model.EmailReqs{Offset: 0})

	// Get expected
	rows, err := suite.db.conn.Query(`
		SELECT emails.id, emails.rcpt, emails.sender, emails.content, emails.received_at, emails.next_retry,
		       emails.exhausted, emails.sent, emails.approved, COALESCE(lists.name, '')
		FROM emails LEFT JOIN lists ON emails.list = lists.id
		WHERE emails.id > $1 ORDER BY id LIMIT $2;
	`, 0, suite.db.batchSize)
	suite.Require().NoError(err)

	var expected []model.Email
	for rows.Next() {
		email := model.Email{}
		err = rows.Scan(&email.Id, pq.Array(&email.Rcpt), &email.Sender, &email.Content, &email.ReceivedAt,
			&email.NextRetry, &email.Exhausted, &email.Sent, &email.Approved, &email.ListName)
		suite.Require().NoError(err)

		expected = append(expected, email)
	}

	suite.Nil(actualErr)
	suite.Equal(suite.db.batchSize, actual.Offset)
	suite.Len(actual.Emails, suite.db.batchSize)
	suite.Equal(expected, actual.Emails)
}

func (suite *DbEmailsSuite) TestGetAllEmailsWithOffset() {
	// Run function
	actual, actualErr := suite.db.GetAllEmails(&model.EmailReqs{Offset: 3})

	// Get expected
	rows, err := suite.db.conn.Query(`
		SELECT emails.id, emails.rcpt, emails.sender, emails.content, emails.received_at, emails.next_retry,
		       emails.exhausted, emails.sent, emails.approved, COALESCE(lists.name, '')
		FROM emails LEFT JOIN lists ON emails.list = lists.id
		WHERE emails.id > $1 ORDER BY id LIMIT $2;
	`, 3, suite.db.batchSize)
	suite.Require().NoError(err)

	var expected []model.Email
	for rows.Next() {
		email := model.Email{}
		err = rows.Scan(&email.Id, pq.Array(&email.Rcpt), &email.Sender, &email.Content, &email.ReceivedAt,
			&email.NextRetry, &email.Exhausted, &email.Sent, &email.Approved, &email.ListName)
		suite.Require().NoError(err)

		expected = append(expected, email)
	}

	suite.Nil(actualErr)
	suite.GreaterOrEqual(suite.db.batchSize+len(actual.Emails), actual.Offset)
	suite.Greater(len(actual.Emails), 1)
	suite.Subset(actual.Emails, expected)
}

func (suite *DbEmailsSuite) TestGetAllEmailsFromList() {
	// Run function
	list := 1
	actual, actualErr := suite.db.GetAllEmails(&model.EmailReqs{Offset: 0, List: &list})

	// Get expected
	rows, err := suite.db.conn.Query(`
		SELECT emails.id, emails.rcpt, emails.sender, emails.content, emails.received_at, emails.next_retry,
		       emails.exhausted, emails.sent, emails.approved, COALESCE(lists.name, '')
		FROM emails LEFT JOIN lists ON emails.list = lists.id
		WHERE emails.id > $1 AND emails.list = $2 ORDER BY id LIMIT $3;
	`, 0, 1, suite.db.batchSize)
	suite.Require().NoError(err)

	var expected []model.Email
	for rows.Next() {
		email := model.Email{}
		err = rows.Scan(&email.Id, pq.Array(&email.Rcpt), &email.Sender, &email.Content, &email.ReceivedAt,
			&email.NextRetry, &email.Exhausted, &email.Sent, &email.Approved, &email.ListName)
		suite.Require().NoError(err)

		expected = append(expected, email)
	}

	suite.Nil(actualErr)
	suite.Equal(len(actual.Emails), actual.Offset)
	suite.Len(actual.Emails, 3)
	suite.Subset(actual.Emails, expected)
}

func (suite *DbEmailsSuite) TestGetAllEmailsWhereArchived() {
	// Run function
	actual, actualErr := suite.db.GetAllEmails(&model.EmailReqs{Offset: 0, Archived: true})

	// Get expected
	rows, err := suite.db.conn.Query(`
		SELECT emails.id, emails.rcpt, emails.sender, emails.content, emails.received_at, emails.next_retry,
		       emails.exhausted, emails.sent, emails.approved, COALESCE(lists.name, '')
		FROM emails LEFT JOIN lists ON emails.list = lists.id
		WHERE emails.id > $1 AND emails.sent = true ORDER BY id LIMIT $2;
	`, 0, suite.db.batchSize)
	suite.Require().NoError(err)

	var expected []model.Email
	for rows.Next() {
		email := model.Email{}
		err = rows.Scan(&email.Id, pq.Array(&email.Rcpt), &email.Sender, &email.Content, &email.ReceivedAt,
			&email.NextRetry, &email.Exhausted, &email.Sent, &email.Approved, &email.ListName)
		suite.Require().NoError(err)

		expected = append(expected, email)
	}

	suite.Nil(actualErr)
	suite.NotEmpty(actual.Emails)
	suite.Subset(actual.Emails, expected)
}

func (suite *DbEmailsSuite) TestGetAllEmailsWhereExhausted() {
	// Run function
	actual, actualErr := suite.db.GetAllEmails(&model.EmailReqs{Offset: 0, Archived: true})

	// Get expected
	rows, err := suite.db.conn.Query(`
		SELECT emails.id, emails.rcpt, emails.sender, emails.content, emails.received_at, emails.next_retry,
		       emails.exhausted, emails.sent, emails.approved, COALESCE(lists.name, '')
		FROM emails LEFT JOIN lists ON emails.list = lists.id
		WHERE emails.id > $1 AND emails.exhausted = 0 ORDER BY id LIMIT $2;
	`, 0, suite.db.batchSize)
	suite.Require().NoError(err)

	var expected []model.Email
	for rows.Next() {
		email := model.Email{}
		err = rows.Scan(&email.Id, pq.Array(&email.Rcpt), &email.Sender, &email.Content, &email.ReceivedAt,
			&email.NextRetry, &email.Exhausted, &email.Sent, &email.Approved, &email.ListName)
		suite.Require().NoError(err)

		expected = append(expected, email)
	}

	suite.Nil(actualErr)
	suite.NotEmpty(actual.Emails)
	suite.Subset(actual.Emails, expected)
}

func (suite *DbEmailsSuite) TestGetAllEmailsWherePendingApproval() {
	// Run function
	actual, actualErr := suite.db.GetAllEmails(&model.EmailReqs{Offset: 0, PendingApproval: true})

	// Get expected
	rows, err := suite.db.conn.Query(`
		SELECT emails.id, emails.rcpt, emails.sender, emails.content, emails.received_at, emails.next_retry,
		       emails.exhausted, emails.sent, emails.approved, COALESCE(lists.name, '')
		FROM emails LEFT JOIN lists ON emails.list = lists.id
		WHERE emails.id > $1 AND emails.approved = false ORDER BY id LIMIT $2;
	`, 0, suite.db.batchSize)
	suite.Require().NoError(err)

	var expected []model.Email
	for rows.Next() {
		email := model.Email{}
		err = rows.Scan(&email.Id, pq.Array(&email.Rcpt), &email.Sender, &email.Content, &email.ReceivedAt,
			&email.NextRetry, &email.Exhausted, &email.Sent, &email.Approved, &email.ListName)
		suite.Require().NoError(err)

		expected = append(expected, email)
	}

	suite.Nil(actualErr)
	suite.NotEmpty(actual.Emails)
	suite.Subset(actual.Emails, expected)
}

func (suite *DbEmailsSuite) TestGetEmail() {
	// Run function
	actual, _ := suite.db.GetEmail(1)

	// Get expected
	expected := &model.Email{}
	err := suite.db.conn.QueryRow(`
		SELECT emails.id, emails.rcpt, emails.sender, emails.content, emails.received_at, emails.next_retry,
		       emails.exhausted, emails.sent, emails.approved, COALESCE(lists.name, '')
		FROM emails LEFT JOIN lists ON emails.list = lists.id WHERE emails.id = $1
	`, 1,
	).Scan(&expected.Id, pq.Array(&expected.Rcpt), &expected.Sender, &expected.Content, &expected.ReceivedAt,
		&expected.NextRetry, &expected.Exhausted, &expected.Sent, &expected.Approved, &expected.ListName)
	suite.Require().NoError(err)

	suite.Equal(expected, actual)
}

func (suite *DbEmailsSuite) TestGetEmailReturnsNoRowsOnInvalidId() {
	_, err := suite.db.GetList(0)

	suite.Equal(ErrNoRows, err.Code)
}

func (suite *DbEmailsSuite) TestGetEmailList() {
	// Run function
	actual, _ := suite.db.GetEmailList(1)

	// Get expected
	expected := 1
	err := suite.db.conn.QueryRow(`
		SELECT list FROM emails WHERE id = $1
	`, 1).Scan(&expected)
	suite.Require().NoError(err)

	suite.Equal(expected, actual)
}

func (suite *DbEmailsSuite) TestGetEmailListReturnsNoRowsOnInvalidId() {
	_, err := suite.db.GetEmailList(0)

	suite.Equal(ErrNoRows, err.Code)
}
