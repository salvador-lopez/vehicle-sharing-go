//go:build unit

package gorm_test

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	gormdb "gorm.io/gorm"
	"gorm.io/gorm/logger"

	databasegorm "vehicle-sharing-go/pkg/database/gorm"
)

type failingDBConnector struct {
	gormdb.ConnPool
	err error
}

type mockSQLPoolPlugin struct {
	pool *sql.DB
}

func (p mockSQLPoolPlugin) Name() string {
	return "mock_sql_pool"
}

func (p mockSQLPoolPlugin) Initialize(db *gormdb.DB) error {
	originalPool, err := db.DB()
	if err != nil {
		return err
	}
	if err := originalPool.Close(); err != nil {
		return err
	}
	db.ConnPool = p.pool
	db.Statement.ConnPool = p.pool
	return nil
}

func (c failingDBConnector) GetDBConn() (*sql.DB, error) {
	return nil, c.err
}

type failingDBConnectorPlugin struct {
	t   *testing.T
	err error
}

func (p failingDBConnectorPlugin) Name() string {
	return "failing_db_connector"
}

func (p failingDBConnectorPlugin) Initialize(db *gormdb.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	p.t.Cleanup(func() {
		require.NoError(p.t, sqlDB.Close())
	})
	db.Statement.ConnPool = failingDBConnector{ConnPool: db.Statement.ConnPool, err: p.err}
	return nil
}

func newConstructorTestConnection(t *testing.T, config *databasegorm.Config, dsnConfig databasegorm.DSNConfig, gormConfig *gormdb.Config) (*databasegorm.Connection, error) {
	t.Helper()

	connection, err := databasegorm.NewConnectionFromConfig(
		config,
		dsnConfig,
		gormConfig,
	)
	if err != nil {
		return connection, err
	}
	require.NotNil(t, connection)
	require.NotNil(t, connection.Db())

	sqlDB, err := connection.Db().DB()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, sqlDB.Close())
	})

	return connection, nil
}

func TestNewConnectionFromConfig(t *testing.T) {
	const maxOpenConnections = 5
	const query = "SELECT 1"
	dbConnectionErr := errors.New("failed to retrieve SQL connection")
	dsnConfig := databasegorm.DSNConfig{
		UserName:     "gorm",
		Password:     "gorm",
		DatabaseName: "gorm",
		Host:         "localhost",
		Port:         9910,
	}
	invalidDSNConfig := dsnConfig
	invalidDSNConfig.DatabaseName = "gorm?timeout=invalid&"

	tests := []struct {
		name                 string
		config               databasegorm.Config
		dsnConfig            databasegorm.DSNConfig
		prepareStatements    bool
		maxOpenConnections   int
		expectedErrorMessage string
		expectedError        error
		setup                func(*testing.T, *gormdb.Config)
		verify               func(*testing.T, *sql.DB)
	}{
		{
			name: "sets maximum open connections",
			config: databasegorm.Config{
				MaxOpenConns:              maxOpenConnections,
				SkipInitializeWithVersion: true,
			},
			maxOpenConnections: maxOpenConnections,
			dsnConfig:          dsnConfig,
		},
		{
			name: "sets maximum idle connections",
			config: databasegorm.Config{
				MaxIdleConns:              1,
				SkipInitializeWithVersion: true,
			},
			dsnConfig: dsnConfig,
			setup: func(t *testing.T, config *gormdb.Config) {
				pool, mock, err := sqlmock.New()
				require.NoError(t, err)
				mock.ExpectClose()
				mock.ExpectClose()
				t.Cleanup(func() {
					require.NoError(t, mock.ExpectationsWereMet())
				})
				plugin := mockSQLPoolPlugin{pool: pool}
				config.Plugins = map[string]gormdb.Plugin{plugin.Name(): plugin}
			},
			verify: func(t *testing.T, pool *sql.DB) {
				first, err := pool.Conn(context.Background())
				require.NoError(t, err)
				second, err := pool.Conn(context.Background())
				require.NoError(t, err)
				require.NoError(t, first.Close())
				require.NoError(t, second.Close())

				require.Equal(t, 1, pool.Stats().Idle)
			},
		},
		{
			name: "allows unlimited open connections for a negative limit",
			config: databasegorm.Config{
				MaxOpenConns:              -1,
				SkipInitializeWithVersion: true,
			},
			maxOpenConnections: 0,
			dsnConfig:          dsnConfig,
		},
		{
			name: "enables query logging",
			config: databasegorm.Config{
				LogQueries:                true,
				SkipInitializeWithVersion: true,
			},
			dsnConfig: dsnConfig,
		},
		{
			name: "preserves prepared statement configuration",
			config: databasegorm.Config{
				SkipInitializeWithVersion: true,
			},
			dsnConfig:         dsnConfig,
			prepareStatements: true,
		},
		{
			name: "returns SQL connection retrieval error",
			config: databasegorm.Config{
				SkipInitializeWithVersion: true,
			},
			dsnConfig:     dsnConfig,
			expectedError: dbConnectionErr,
			setup: func(t *testing.T, config *gormdb.Config) {
				plugin := failingDBConnectorPlugin{t: t, err: dbConnectionErr}
				config.Plugins = map[string]gormdb.Plugin{plugin.Name(): plugin}
			},
		},
		{
			name: "returns GORM opening error for invalid DSN",
			config: databasegorm.Config{
				SkipInitializeWithVersion: true,
			},
			dsnConfig:            invalidDSNConfig,
			expectedErrorMessage: `time: invalid duration "invalid"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			gormConfig := &gormdb.Config{
				DryRun:               true,
				DisableAutomaticPing: true,
				PrepareStmt:          tt.prepareStatements,
				Logger: logger.New(log.New(&output, "", 0), logger.Config{
					LogLevel: logger.Silent,
				}),
			}
			if tt.setup != nil {
				tt.setup(t, gormConfig)
			}
			connection, err := newConstructorTestConnection(t, &tt.config, tt.dsnConfig, gormConfig)
			if tt.expectedError != nil {
				require.ErrorIs(t, err, tt.expectedError)
				require.Nil(t, connection)
				return
			}
			if tt.expectedErrorMessage != "" {
				require.ErrorContains(t, err, tt.expectedErrorMessage)
				require.Nil(t, connection)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.prepareStatements, connection.Db().PrepareStmt)
			sqlDB, err := connection.Db().DB()
			require.NoError(t, err)
			require.Equal(t, tt.maxOpenConnections, sqlDB.Stats().MaxOpenConnections)

			if tt.verify != nil {
				tt.verify(t, sqlDB)
			}

			require.NoError(t, connection.Db().WithContext(context.Background()).Exec(query).Error)
			if tt.config.LogQueries {
				require.Contains(t, output.String(), query)
			} else {
				require.Empty(t, output.String())
			}
		})
	}
}

func TestConnectionIsDuplicateEntryErr(t *testing.T) {
	connection := &databasegorm.Connection{}
	tests := []struct {
		name        string
		err         error
		isDuplicate bool
	}{
		{
			name: "true",
			err: &mysql.MySQLError{
				Number:  1062,
				Message: "Duplicate entry",
			},
			isDuplicate: true,
		},
		{
			name:        "false",
			err:         errors.New("database error"),
			isDuplicate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.isDuplicate, connection.IsDuplicateEntryErr(tt.err))
		})
	}
}
