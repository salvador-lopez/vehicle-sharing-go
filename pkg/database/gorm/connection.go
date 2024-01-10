package gorm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	UserName        string
	Password        string
	DatabaseName    string
	Host            string
	Port            int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	MaxOpenConns    int
	Logger          *log.Logger
	LogQueries      bool
}

type Connection struct {
	db     *gorm.DB
	tx     *gorm.DB
	logger log.Logger
}

func (c *Connection) Db() *gorm.DB {
	if c.tx != nil {
		return c.tx
	}

	return c.db
}

func NewConnectionFromConfig(c *Config) (*Connection, error) {
	dsn := fmt.Sprintf(`%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=UTC`,
		c.UserName,
		c.Password,
		c.Host,
		c.Port,
		c.DatabaseName,
	)

	db, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{
		Logger: logger.Default,
	})

	if err != nil {
		return nil, err
	}

	if c.LogQueries {
		db.Logger = db.Debug().Logger
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(c.ConnMaxLifetime)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)

	return &Connection{db: db}, nil
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
