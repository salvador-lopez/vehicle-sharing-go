package gorm

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
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
	logger log.Logger
}

func (c *Connection) Db() *gorm.DB {
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

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
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

func (c *Connection) Close() error {
	if c == nil || c.db == nil {
		return nil
	}

	sqlDb, err := c.db.DB()
	if err != nil {
		return err
	}

	return sqlDb.Close()
}

func (c *Connection) Transaction(ctx context.Context, f func(context.Context) error) error {
	return c.db.Transaction(func(tx *gorm.DB) error {
		c.db = tx

		return f(ctx)
	})
}
