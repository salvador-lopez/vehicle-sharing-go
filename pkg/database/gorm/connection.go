package gorm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	MaxIdleConns              int
	ConnMaxLifetime           time.Duration
	MaxOpenConns              int
	LogQueries                bool
	SkipInitializeWithVersion bool
}

type DSNConfig struct {
	UserName     string
	Password     string
	DatabaseName string
	Host         string
	Port         int
}

type Connection struct {
	db *gorm.DB
	tx *gorm.DB
}

func buildDSN(config DSNConfig) string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		config.UserName,
		config.Password,
		config.Host,
		config.Port,
		config.DatabaseName,
	)
}

func NewConnectionFromConfig(connectionConfig *Config, dsnConfig DSNConfig, gormConfig *gorm.Config) (*Connection, error) {
	db, err := gorm.Open(gormmysql.New(gormmysql.Config{
		DSN:                       buildDSN(dsnConfig),
		SkipInitializeWithVersion: connectionConfig.SkipInitializeWithVersion,
	}), gormConfig)
	if err != nil {
		return nil, err
	}

	if connectionConfig.LogQueries {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(connectionConfig.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(connectionConfig.ConnMaxLifetime)
	sqlDB.SetMaxOpenConns(connectionConfig.MaxOpenConns)

	return &Connection{db: db}, nil
}

func (c *Connection) Db() *gorm.DB {
	if c.tx != nil {
		return c.tx
	}

	return c.db
}

func (c *Connection) Transaction(ctx context.Context, f func(context.Context) error) error {
	return c.db.Transaction(func(tx *gorm.DB) error {
		defer func() { c.tx = nil }()
		c.tx = tx

		return f(ctx)
	})
}

func (c *Connection) IsDuplicateEntryErr(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}

	return false
}
