package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
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
	port := util.MockConfigs["SQL_PORT"]
	c, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase("goformail"),
		postgres.WithUsername("goformail"),
		postgres.WithPassword("password"),
		postgres.WithInitScripts("scripts/lists_test.sql"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp")),
		testcontainers.WithHostConfigModifier(func(hostConfig *container.HostConfig) {
			hostConfig.AutoRemove = true
			hostConfig.PortBindings = map[nat.Port][]nat.PortBinding{"5432/tcp": {{HostPort: port}}}
			hostConfig.Tmpfs = map[string]string{"/var/lib/postgresql/data": "rw"}
		}),
	)
	suite.NoError(err)

	info := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		util.MockConfigs["SQL_ADDRESS"], port, util.MockConfigs["SQL_USER"], util.MockConfigs["SQL_PASSWORD"],
		util.MockConfigs["SQL_DB_NAME"])
	db, err := sql.Open("postgres", info)
	suite.NoError(err)

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
	var expected model.List
	err := suite.db.conn.QueryRow(`
		SELECT * FROM lists WHERE name = $1
	`, "get-test-0").Scan(new(int), &expected.Name, pq.Array(&expected.Recipients))
	suite.NoError(err)

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
	var actual model.List
	err := suite.db.conn.QueryRow(`
		SELECT * FROM lists WHERE id = $1
	`, id).Scan(new(int), &actual.Name, pq.Array(&actual.Recipients))
	suite.NoError(err)

	suite.Equal(expected, &actual)
}

func (suite *DbListsSuite) TestCreateListReturnsErrorOnDuplicate() {
	list := createList("create-test-1")
	suite.db.CreateList(list)
	_, err := suite.db.CreateList(list)

	suite.Equal(ErrDuplicate, err.Code)
}

func createList(name string) *model.List {
	return &model.List{Name: name, Recipients: []string{"example@domain.tld"}}
}
