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
)

type DbUsersSuite struct {
	suite.Suite
	c  *postgres.PostgresContainer
	db *Db
}

func TestDbUsersSuite(t *testing.T) {
	suite.Run(t, new(DbUsersSuite))
}

func (suite *DbUsersSuite) SetupSuite() {
	ctx := context.Background()
	c, err := postgres.Run(
		ctx,
		"postgres:latest",
		postgres.WithDatabase("goformail"),
		postgres.WithUsername("goformail"),
		postgres.WithPassword("password"),
		postgres.WithInitScripts("scripts/users_test.sql"),
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

	suite.c, suite.db = c, &Db{conn: db}
}

func (suite *DbUsersSuite) TearDownSuite() {
	if err := testcontainers.TerminateContainer(suite.c); err != nil {
		suite.T().Log(err)
	}
}

func (suite *DbUsersSuite) TestGetUser() {
	// Run function
	actual, _ := suite.db.GetUser(1)

	// Get expected
	expected := model.UserResponse{Id: 1}
	err := suite.db.conn.QueryRow(`
		SELECT email, permissions FROM users WHERE id = $1
	`, 1).Scan(&expected.Email, pq.Array(&expected.Permissions))
	suite.Require().NoError(err)

	suite.Equal(&expected, actual)
}

func (suite *DbUsersSuite) TestGetUserReturnsNoRowsOnInvalidId() {
	_, err := suite.db.GetUser(0)

	suite.Equal(ErrNoRows, err.Code)
}

func (suite *DbUsersSuite) TestCreateUser() {
	// Run function
	expected := createUser("create@test-0.tld")
	id, _ := suite.db.CreateUser(expected, "hash")

	// Check list was created properly
	var actual model.UserRequest
	var actualHash string
	err := suite.db.conn.QueryRow(`
		SELECT email, hash, permissions FROM users WHERE id = $1
	`, id).Scan(&actual.Email, &actualHash, pq.Array(&actual.Permissions))
	suite.Require().NoError(err)

	suite.Equal(expected, &actual)
	suite.Equal("hash", actualHash)
}

func (suite *DbUsersSuite) TestCreateUserReturnsErrorOnDuplicate() {
	user := createUser("create@test-1.tld")
	suite.db.CreateUser(user, "hash")
	_, err := suite.db.CreateUser(user, "hash1")

	suite.Equal(ErrDuplicate, err.Code)
}

func (suite *DbUsersSuite) TestUpdateUser() {
	// Run function
	expected := createUser("update@test-8.tld")
	suite.db.UpdateUser(2, expected, false)

	// Check list was updated properly
	var actual model.UserRequest
	err := suite.db.conn.QueryRow(`
		SELECT email, permissions FROM users WHERE id = $1
	`, 2).Scan(&actual.Email, pq.Array(&actual.Permissions))
	suite.Require().NoError(err)

	suite.Equal(expected, &actual)
}

func (suite *DbUsersSuite) TestUpdateUserUpdatesPartially() {
	// Run function
	expected := &model.UserRequest{Permissions: []string{"ADMIN", "CRT_LIST"}}
	suite.db.UpdateUser(3, expected, false)

	// Check list was patched properly
	var actual model.UserResponse
	err := suite.db.conn.QueryRow(`
		SELECT email, permissions FROM users WHERE id = $1
	`, 3).Scan(&actual.Email, pq.Array(&actual.Permissions))
	suite.Require().NoError(err)

	suite.Equal(&model.UserResponse{Email: "update@test-1.tld", Permissions: expected.Permissions}, &actual)
}

func (suite *DbUsersSuite) TestUpdateOverridePerms() {
	// Run function
	expected := &model.UserRequest{Permissions: []string{}}
	suite.db.UpdateUser(3, expected, true)

	// Check list was updated properly
	var actual model.UserResponse
	err := suite.db.conn.QueryRow(`
		SELECT email, permissions FROM users WHERE id = $1
	`, 3).Scan(&actual.Email, pq.Array(&actual.Permissions))
	suite.Require().NoError(err)

	suite.Equal(&model.UserResponse{Email: "update@test-1.tld", Permissions: expected.Permissions}, &actual)
}

func (suite *DbUsersSuite) TestUpdateUserReturnsNoRowsOnInvalidId() {
	err := suite.db.UpdateUser(0, createUser("update@test-0.tld"), false)

	suite.Equal(ErrNoRows, err.Code)
}

func (suite *DbUsersSuite) TestUpdateUserReturnsDuplicateOnExistingList() {
	err := suite.db.UpdateUser(2, createUser("update@test-1.tld"), false)

	suite.Equal(ErrDuplicate, err.Code)
}

func (suite *DbUsersSuite) TestDeleteUser() {
	// Run function
	suite.db.DeleteUser(4)

	// Check list was patched properly
	err := suite.db.conn.QueryRow(`
		SELECT * FROM users WHERE id = $1
	`, 4).Scan()

	suite.Equal(sql.ErrNoRows, err)
}

func (suite *DbUsersSuite) TestDeleteUserReturnsNoRowsOnInvalidId() {
	err := suite.db.DeleteUser(0)

	suite.Equal(ErrNoRows, err.Code)
}

func (suite *DbUsersSuite) TestGetAllUsers() {
	// Run function
	actual, _ := suite.db.GetAllUsers()

	// Get expected
	var expected []*model.UserResponse
	rows, err := suite.db.conn.Query(`
		SELECT id, email, permissions FROM users
	`)
	suite.Require().NoError(err)
	for rows.Next() {
		user := model.UserResponse{}
		err := rows.Scan(&user.Id, &user.Email, pq.Array(&user.Permissions))
		suite.Require().NoError(err)

		expected = append(expected, &user)
	}

	suite.Equal(&expected, actual)
}

func (suite *DbUsersSuite) TestGetUserPassword() {
	actualId, actualHash, err := suite.db.GetUserPassword("get@test-0.tld")

	suite.Nil(err)
	suite.Equal(1, actualId)
	suite.Equal("hash", actualHash)
}

func (suite *DbUsersSuite) TestGetUserPasswordReturnsNoRowsOnInvalidEmail() {
	_, _, err := suite.db.GetUserPassword("invalid")

	suite.Equal(ErrNoRows, err.Code)
}

func (suite *DbUsersSuite) TestUserExists() {
	actual, err := suite.db.UserExists(1)

	suite.Nil(err)
	suite.True(actual)
}

func (suite *DbUsersSuite) TestUserExistsReturnsFalseOnInvalidId() {
	actual, err := suite.db.UserExists(0)

	suite.Nil(err)
	suite.False(actual)
}

func (suite *DbUsersSuite) TestUsersExist() {
	expected := []int64{1, 2}
	actual, err := suite.db.UsersExist(expected)

	suite.Nil(err)
	suite.Equal(expected, actual)
}

func (suite *DbUsersSuite) TestUsersExistReturnsMissingOnInvalidId() {
	actual, err := suite.db.UsersExist([]int64{1, 2, 999})

	suite.Nil(err)
	suite.Equal([]int64{1, 2}, actual)
}

func (suite *DbUsersSuite) TestUsersExistReturnsMissingOnAllInvalidId() {
	actual, err := suite.db.UsersExist([]int64{0, 999})

	suite.Nil(err)
	suite.Empty(actual)
}

func (suite *DbUsersSuite) TestGetUserPerms() {
	// Run function
	actual, _ := suite.db.GetUserPerms(1)

	// Get expected
	var expected []string
	err := suite.db.conn.QueryRow(`
		SELECT permissions FROM users WHERE id = $1
	`, 1).Scan(pq.Array(&expected))
	suite.Require().NoError(err)

	suite.Equal(expected, actual)
}

func (suite *DbUsersSuite) TestGetUserPermsReturnsNoRowsOnInvalidId() {
	_, err := suite.db.GetUserPerms(0)

	suite.Equal(ErrNoRows, err.Code)
}

func (suite *DbUsersSuite) TestGetUserPermsAndModStatus() {
	// Run function
	actualPerms, actualStatus, _ := suite.db.GetUserPermsAndModStatus(1, 1)

	// Get expected
	var expected []string
	err := suite.db.conn.QueryRow(`
		SELECT permissions FROM users WHERE id = $1
	`, 1).Scan(pq.Array(&expected))
	suite.Require().NoError(err)

	suite.Equal(expected, actualPerms)
	suite.True(actualStatus)
}

func (suite *DbUsersSuite) TestGetUserPermsAndModStatusWhenNotMod() {
	// Run function
	actualPerms, actualStatus, _ := suite.db.GetUserPermsAndModStatus(2, 1)

	// Get expected
	var expected []string
	err := suite.db.conn.QueryRow(`
		SELECT permissions FROM users WHERE id = $1
	`, 1).Scan(pq.Array(&expected))
	suite.Require().NoError(err)

	suite.Equal(expected, actualPerms)
	suite.False(actualStatus)
}

func (suite *DbUsersSuite) TestGetUserPermsAndModStatusWhenNoList() {
	// Run function
	actualPerms, actualStatus, _ := suite.db.GetUserPermsAndModStatus(1, 0)

	// Get expected
	var expected []string
	err := suite.db.conn.QueryRow(`
		SELECT permissions FROM users WHERE id = $1
	`, 1).Scan(pq.Array(&expected))
	suite.Require().NoError(err)

	suite.Equal(expected, actualPerms)
	suite.False(actualStatus)
}

func (suite *DbUsersSuite) TestTestGetUserPermsAndModStatusReturnsNoRowsOnInvalidId() {
	_, _, err := suite.db.GetUserPermsAndModStatus(0, 1)

	suite.Equal(ErrNoRows, err.Code)
}

func createUser(email string) *model.UserRequest {
	return &model.UserRequest{Email: email, Permissions: []string{"ADMIN", "CRT_LIST"}}
}
