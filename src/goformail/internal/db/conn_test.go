package db

import (
	"context"
	"database/sql"
	"github.com/docker/docker/api/types/container"
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

var expectedPropColumns = []column{
	{"id", "integer", "NO", sql.NullString{String: "nextval('props_id_seq'::regclass)", Valid: true}},
	{name: "jwt_secret", ctype: "bytea", nullable: "YES", dflt: sql.NullString{String: "", Valid: false}},
}

var expectedListColumns = []column{
	{"id", "integer", "NO", sql.NullString{String: "nextval('lists_id_seq'::regclass)", Valid: true}},
	{name: "mods", ctype: "ARRAY", nullable: "YES", dflt: sql.NullString{String: "", Valid: false}},
	{name: "locked", ctype: "boolean", nullable: "YES", dflt: sql.NullString{String: "false", Valid: true}},
	{name: "name", ctype: "text", nullable: "NO", dflt: sql.NullString{String: "", Valid: false}},
	{name: "recipients", ctype: "ARRAY", nullable: "YES", dflt: sql.NullString{String: "", Valid: false}},
	{name: "approved_senders", ctype: "ARRAY", nullable: "YES", dflt: sql.NullString{String: "", Valid: false}},
}

var expectedUserColumns = []column{
	{"id", "integer", "NO", sql.NullString{String: "nextval('users_id_seq'::regclass)", Valid: true}},
	{name: "email", ctype: "text", nullable: "NO", dflt: sql.NullString{String: "", Valid: false}},
	{name: "hash", ctype: "text", nullable: "NO", dflt: sql.NullString{String: "", Valid: false}},
	{name: "permissions", ctype: "ARRAY", nullable: "YES", dflt: sql.NullString{String: "", Valid: false}},
}

var expectedEmailColumns = []column{
	{name: "approved", ctype: "boolean", nullable: "YES", dflt: sql.NullString{String: "false", Valid: true}},
	{name: "sent", ctype: "boolean", nullable: "YES", dflt: sql.NullString{String: "false", Valid: true}},
	{name: "list", ctype: "integer", nullable: "YES", dflt: sql.NullString{String: "", Valid: false}},
	{"id", "integer", "NO", sql.NullString{String: "nextval('emails_id_seq'::regclass)", Valid: true}},
	{"received_at", "timestamp without time zone", "NO", sql.NullString{String: "", Valid: false}},
	{"next_retry", "timestamp without time zone", "YES", sql.NullString{String: "", Valid: false}},
	{name: "exhausted", ctype: "integer", nullable: "YES", dflt: sql.NullString{String: "3", Valid: true}},
	{name: "rcpt", ctype: "ARRAY", nullable: "NO", dflt: sql.NullString{String: "", Valid: false}},
	{name: "sender", ctype: "text", nullable: "NO", dflt: sql.NullString{String: "", Valid: false}},
	{name: "content", ctype: "text", nullable: "NO", dflt: sql.NullString{String: "", Valid: false}},
}

func TestInitDBCreatesTablesAndJWTCreatesThenFetchesSecret(t *testing.T) {
	// Define the postgres container config
	ctx := context.Background()
	c, err := postgres.Run(
		ctx,
		"postgres:latest",
		postgres.WithDatabase("goformail"),
		postgres.WithUsername("goformail"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp")),
		testcontainers.WithHostConfigModifier(func(hostConfig *container.HostConfig) {
			hostConfig.AutoRemove = true
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
	host, _ := c.Host(ctx)
	mappedPort, _ := c.MappedPort(ctx, "5432")
	containerConfig := map[string]string{
		"SQL_ADDRESS":  host,
		"SQL_PORT":     mappedPort.Port(),
		"SQL_USER":     util.MockConfigs["SQL_USER"],
		"SQL_PASSWORD": util.MockConfigs["SQL_PASSWORD"],
		"SQL_DB_NAME":  util.MockConfigs["SQL_DB_NAME"],
	}
	db := InitDB(containerConfig)

	// Check correct props columns are present
	rows, err := db.conn.Query(`
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_name = 'props';
	`)
	require.NoError(t, err)
	var columns []column
	for rows.Next() {
		column := column{}
		err := rows.Scan(&column.name, &column.ctype, &column.nullable, &column.dflt)
		require.NoError(t, err)

		columns = append(columns, column)
	}

	assert.Equal(t, expectedPropColumns, columns)

	// Check correct lists columns are present
	rows, err = db.conn.Query(`
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_name = 'lists';
	`)
	require.NoError(t, err)
	columns = []column{}
	for rows.Next() {
		column := column{}
		err := rows.Scan(&column.name, &column.ctype, &column.nullable, &column.dflt)
		require.NoError(t, err)

		columns = append(columns, column)
	}

	assert.Equal(t, expectedListColumns, columns)

	// Check correct users columns are present
	rows, err = db.conn.Query(`
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_name = 'users';
	`)
	require.NoError(t, err)
	columns = []column{}
	for rows.Next() {
		column := column{}
		err := rows.Scan(&column.name, &column.ctype, &column.nullable, &column.dflt)
		require.NoError(t, err)

		columns = append(columns, column)
	}

	assert.Equal(t, expectedUserColumns, columns)

	// Check correct users columns are present
	rows, err = db.conn.Query(`
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_name = 'emails';
	`)
	require.NoError(t, err)
	columns = []column{}
	for rows.Next() {
		column := column{}
		err := rows.Scan(&column.name, &column.ctype, &column.nullable, &column.dflt)
		require.NoError(t, err)

		columns = append(columns, column)
	}

	assert.Equal(t, expectedEmailColumns, columns)

	secretInDb := make([]byte, 32)
	err = db.conn.QueryRow("SELECT jwt_secret FROM props LIMIT 1;").Scan(&secretInDb)
	require.NoError(t, err)
	require.Nil(t, secretInDb)

	secretGenerated := db.GetJwtSecret()
	err = db.conn.QueryRow("SELECT jwt_secret FROM props LIMIT 1;").Scan(&secretInDb)
	require.NoError(t, err)
	require.NotNil(t, *secretGenerated)
	require.Equal(t, *secretGenerated, secretInDb)

	secretGenerated = db.GetJwtSecret()
	require.NoError(t, err)
	require.Equal(t, *secretGenerated, secretInDb)
}
