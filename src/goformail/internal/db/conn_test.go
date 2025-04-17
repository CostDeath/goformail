package db

import (
	"context"
	"database/sql"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"testing"
)

type column struct {
	name, ctype, nullable string
	dflt                  sql.NullString
}

var expectedColumns = []column{
	{"id", "integer", "NO", sql.NullString{String: "nextval('lists_id_seq'::regclass)", Valid: true}},
	{name: "name", ctype: "text", nullable: "NO", dflt: sql.NullString{String: "", Valid: false}},
	{name: "recipients", ctype: "ARRAY", nullable: "YES", dflt: sql.NullString{String: "", Valid: false}},
}

func TestInitDBCreatesTables(t *testing.T) {
	// Define the postgres container config
	port := util.MockConfigs["SQL_PORT"]
	c, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase("goformail"),
		postgres.WithUsername("goformail"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp")),
		testcontainers.WithHostConfigModifier(func(hostConfig *container.HostConfig) {
			hostConfig.AutoRemove = true
			hostConfig.PortBindings = map[nat.Port][]nat.PortBinding{"5432/tcp": {{HostPort: port}}}
			hostConfig.Tmpfs = map[string]string{"/var/lib/postgresql/data": "rw"}
		}),
	)
	require.NoError(t, err)

	defer func() {
		if err := testcontainers.TerminateContainer(c); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	}()

	// Run function
	db := InitDB(util.MockConfigs)

	// Check correct tables are present
	rows, err := db.conn.Query(`
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_name = 'lists';
	`)
	require.NoError(t, err)
	var columns []column
	for rows.Next() {
		column := column{}
		err := rows.Scan(&column.name, &column.ctype, &column.nullable, &column.dflt)
		require.NoError(t, err)

		columns = append(columns, column)
	}

	assert.Equal(t, expectedColumns, columns)
}
