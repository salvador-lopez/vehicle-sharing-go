//go:build integration

package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// mysqlHelper manages a raw *sql.DB connection used to set up and tear down
// the table fixtures that the binlog Source is expected to stream changes
// from. It is intentionally independent from the production code under
// test: it exists purely to drive real MySQL statements (DDL/DML) as an
// external actor would.
type mysqlHelper struct {
	db *sql.DB
}

// newMySQLHelper opens a connection to the dedicated integration test
// database, configured via the MYSQL_HOST/PORT/USER/PASSWORD/
// MYSQL_MESSAGE_RELAY_DATABASE env vars (falling back to the docker-compose
// defaults). The database must be independent from any business context's
// database (e.g. "inventory"), since this package is reusable
// infrastructure and must not be coupled to a specific context; its
// default ("message_relay_test") must match the database provisioned for
// this user in init.sql.
func newMySQLHelper() (*mysqlHelper, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		envOrDefault("MYSQL_USER", "inventory"),
		envOrDefault("MYSQL_PASSWORD", "inventory"),
		envOrDefault("MYSQL_HOST", "localhost"),
		envOrDefault("MYSQL_PORT", "3308"),
		envOrDefault("MYSQL_MESSAGE_RELAY_DATABASE", "message_relay_test"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &mysqlHelper{db: db}, nil
}

// recreateTable drops (if present) and creates a fresh table with a single
// id/value schema, so each test starts from a clean, isolated state.
func (h *mysqlHelper) recreateTable(name string) error {
	if err := h.dropTable(name); err != nil {
		return err
	}

	_, err := h.db.Exec(fmt.Sprintf(
		"CREATE TABLE %s (id INT PRIMARY KEY, value VARCHAR(255) NOT NULL)",
		name,
	))
	return err
}

func (h *mysqlHelper) dropTable(name string) error {
	_, err := h.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", name))
	return err
}

func (h *mysqlHelper) insertRow(table string, id int, value string) error {
	_, err := h.db.Exec(fmt.Sprintf("INSERT INTO %s (id, value) VALUES (?, ?)", table), id, value)
	return err
}

func (h *mysqlHelper) close() error {
	return h.db.Close()
}
