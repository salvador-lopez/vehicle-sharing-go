//go:build unit

package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"vehicle-sharing-go/pkg/database/gorm/model"
)

func TestOutboxRecordTableName(t *testing.T) {
	record := &model.OutboxRecord{}

	require.Equal(t, "outbox", record.TableName())
}
