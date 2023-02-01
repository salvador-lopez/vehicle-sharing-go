package gorm

import (
	"time"

	"github.com/stretchr/testify/require"
)

func requireEqualDates(expected time.Time, actual time.Time, require *require.Assertions) {
	require.Equal(expected.UTC().Format(time.RFC3339), actual.UTC().Format(time.RFC3339))
}
