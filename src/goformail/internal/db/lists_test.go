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

type DbListsSuite struct {
	suite.Suite
	c  *postgres.PostgresContainer
	db *Db
}

func TestDbListsSuite(t *testing.T) {
	suite.Run(t, new(DbListsSuite))
}

func (suite *DbListsSuite) SetupSuite() {
	ctx := context.Background()
	c, err := postgres.Run(
		ctx,
		"postgres:latest",
		postgres.WithDatabase("goformail"),
		postgres.WithUsername("goformail"),
		postgres.WithPassword("password"),
		postgres.WithInitScripts("scripts/lists_test.sql"),
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

func (suite *DbListsSuite) TearDownSuite() {
	if err := testcontainers.TerminateContainer(suite.c); err != nil {
		suite.T().Log(err)
	}
}

func (suite *DbListsSuite) TestGetList() {
	// Run function
	actual, _ := suite.db.GetList(1)

	// Get expected
	expected := model.ListResponse{Id: 1}
	err := suite.db.conn.QueryRow(`
		SELECT name, recipients, mods, approved_senders FROM lists WHERE name = $1
	`, "get-test-0").Scan(&expected.Name, pq.Array(&expected.Recipients), pq.Array(&expected.Mods),
		pq.Array(&expected.ApprovedSenders))
	suite.Require().NoError(err)

	suite.Equal(&expected, actual)
}

func (suite *DbListsSuite) TestGetListReturnsNoRowsOnInvalidId() {
	_, err := suite.db.GetList(0)

	suite.Equal(ErrNoRows, err.Code)
}

func (suite *DbListsSuite) TestCreateList() {
	// Run function
	expected := createList("create-test-0")
	id, _ := suite.db.CreateList(expected)

	// Check list was created properly
	var actual model.ListRequest
	err := suite.db.conn.QueryRow(`
		SELECT name, recipients, mods, approved_senders FROM lists WHERE id = $1
	`, id).Scan(&actual.Name, pq.Array(&actual.Recipients), pq.Array(&actual.Mods), pq.Array(&actual.ApprovedSenders))
	suite.Require().NoError(err)

	suite.Equal(expected, &actual)
}

func (suite *DbListsSuite) TestCreateListReturnsErrorOnDuplicate() {
	list := createList("create-test-1")
	suite.db.CreateList(list)
	_, err := suite.db.CreateList(list)

	suite.Equal(ErrDuplicate, err.Code)
}

func (suite *DbListsSuite) TestPatchList() {
	// Run function
	expected := createList("patch-test-8")
	suite.db.PatchList(2, expected, nil)

	// Check list was patched properly
	var actual model.ListRequest
	err := suite.db.conn.QueryRow(`
		SELECT name, recipients, mods, approved_senders FROM lists WHERE id = $1
	`, 2).Scan(&actual.Name, pq.Array(&actual.Recipients), pq.Array(&actual.Mods), pq.Array(&actual.ApprovedSenders))
	suite.Require().NoError(err)

	suite.Equal(expected, &actual)
}

func (suite *DbListsSuite) TestPatchListUpdatesPartially() {
	// Run function
	expected := &model.ListRequest{Recipients: []string{"example2@domain.tld"}}
	suite.db.PatchList(3, expected, nil)

	// Check list was patched properly
	var actual model.ListRequest
	err := suite.db.conn.QueryRow(`
		SELECT name, recipients, mods, approved_senders FROM lists WHERE id = $1
	`, 3).Scan(&actual.Name, pq.Array(&actual.Recipients), pq.Array(&actual.Mods), pq.Array(&actual.ApprovedSenders))
	suite.Require().NoError(err)

	suite.Equal(&model.ListRequest{Name: "patch-test-1", Recipients: expected.Recipients, Mods: []int64{1},
		ApprovedSenders: []string{"example@domain.tld"}}, &actual)
}

func (suite *DbListsSuite) TestPatchListOverrideRecipients() {
	// Run function
	expected := &model.ListRequest{Recipients: []string{}}
	suite.db.PatchList(4, expected, &model.ListOverrides{Recipients: true})

	// Check list was patched properly
	var actual model.ListRequest
	err := suite.db.conn.QueryRow(`
		SELECT name, recipients FROM lists WHERE id = $1
	`, 4).Scan(&actual.Name, pq.Array(&actual.Recipients))
	suite.Require().NoError(err)

	suite.Equal(&model.ListRequest{Name: "patch-test-2", Recipients: expected.Recipients}, &actual)
}

func (suite *DbListsSuite) TestPatchListOverrideMods() {
	// Run function
	expected := &model.ListRequest{Mods: []int64{}}
	suite.db.PatchList(4, expected, &model.ListOverrides{Mods: true})

	// Check list was patched properly
	var actual model.ListRequest
	err := suite.db.conn.QueryRow(`
		SELECT name, mods FROM lists WHERE id = $1
	`, 4).Scan(&actual.Name, pq.Array(&actual.Mods))
	suite.Require().NoError(err)

	suite.Equal(&model.ListRequest{Name: "patch-test-2", Mods: expected.Mods}, &actual)
}

func (suite *DbListsSuite) TestPatchListOverrideApprovedSenders() {
	// Run function
	expected := &model.ListRequest{ApprovedSenders: []string{}}
	suite.db.PatchList(4, expected, &model.ListOverrides{ApprovedSenders: true})

	// Check list was patched properly
	var actual model.ListRequest
	err := suite.db.conn.QueryRow(`
		SELECT name, approved_senders FROM lists WHERE id = $1
	`, 4).Scan(&actual.Name, pq.Array(&actual.ApprovedSenders))
	suite.Require().NoError(err)

	suite.Equal(&model.ListRequest{Name: "patch-test-2", ApprovedSenders: expected.ApprovedSenders}, &actual)
}

func (suite *DbListsSuite) TestPatchListOverrideAll() {
	// Run function
	expected := &model.ListRequest{Recipients: []string{}}
	suite.db.PatchList(5, expected, &model.ListOverrides{Recipients: true, Mods: true, ApprovedSenders: true})

	// Check list was patched properly
	var actual model.ListRequest
	err := suite.db.conn.QueryRow(`
		SELECT name, recipients, mods, approved_senders FROM lists WHERE id = $1
	`, 5).Scan(&actual.Name, pq.Array(&actual.Recipients), pq.Array(&actual.Mods),
		pq.Array(&actual.ApprovedSenders))
	suite.Require().NoError(err)

	suite.Equal(&model.ListRequest{Name: "patch-test-3", Recipients: expected.Recipients, Mods: expected.Mods,
		ApprovedSenders: expected.ApprovedSenders}, &actual)
}

func (suite *DbListsSuite) TestPatchListReturnsNoRowsOnInvalidId() {
	err := suite.db.PatchList(0, createList("patch-test-0"), nil)

	suite.Equal(ErrNoRows, err.Code)
}

func (suite *DbListsSuite) TestPatchListReturnsDuplicateOnExistingList() {
	err := suite.db.PatchList(2, createList("patch-test-1"), nil)

	suite.Equal(ErrDuplicate, err.Code)
}

func (suite *DbListsSuite) TestDeleteList() {
	// Run function
	suite.db.DeleteList(6)

	// Check list was patched properly
	err := suite.db.conn.QueryRow(`
		SELECT * FROM lists WHERE id = $1
	`, 6).Scan()

	suite.Equal(sql.ErrNoRows, err)
}

func (suite *DbListsSuite) TestDeleteListReturnsNoRowsOnInvalidId() {
	err := suite.db.DeleteList(0)

	suite.Equal(ErrNoRows, err.Code)
}

func (suite *DbListsSuite) TestGetAllList() {
	// Run function
	actual, _ := suite.db.GetAllLists()

	// Get expected
	var expected []*model.ListResponse
	rows, err := suite.db.conn.Query(`
		SELECT * FROM lists
	`)
	suite.Require().NoError(err)
	for rows.Next() {
		list := model.ListResponse{}
		err := rows.Scan(&list.Id, &list.Name, pq.Array(&list.Recipients), pq.Array(&list.Mods),
			pq.Array(&list.ApprovedSenders))
		suite.Require().NoError(err)

		expected = append(expected, &list)
	}

	suite.Equal(&expected, actual)
}

func createList(name string) *model.ListRequest {
	return &model.ListRequest{
		Name:            name,
		Recipients:      []string{"example@domain.tld"},
		Mods:            []int64{1},
		ApprovedSenders: []string{"example2@domain.tld"},
	}
}
