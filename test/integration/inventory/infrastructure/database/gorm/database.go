package gorm

import (
	"fmt"
	"os"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func createDb(require *require.Assertions) *gorm.DB {
	dsn := fmt.Sprintf(`%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC`,
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	require.NoError(err)

	return db
}
