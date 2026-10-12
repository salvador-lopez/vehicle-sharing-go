//go:build unit

package gorm_test

import databasegorm "vehicle-sharing-go/pkg/database/gorm"

func newTestDSNConfig() databasegorm.DSNConfig {
	return databasegorm.DSNConfig{
		UserName:     "gorm",
		Password:     "gorm",
		DatabaseName: "gorm",
		Host:         "localhost",
		Port:         9910,
	}
}
