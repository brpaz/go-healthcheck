package dbcheck_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/brpaz/go-healthcheck/v2/checks"
	"github.com/brpaz/go-healthcheck/v2/checks/dbcheck"
)

// MockDatabaseStatsProvider is a mock implementation of the database stats interface
type MockDatabaseStatsProvider struct {
	mock.Mock
}

func (m *MockDatabaseStatsProvider) Stats() sql.DBStats {
	args := m.Called()
	return args.Get(0).(sql.DBStats)
}

func TestConnectionsCheck_Run(t *testing.T) {
	t.Parallel()

	t.Run("check passes when connections are below maximum", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabaseStatsProvider{}
		stats := sql.DBStats{
			OpenConnections:    10,
			MaxOpenConnections: 100,
		}
		mockDB.On("Stats").Return(stats)

		check := dbcheck.NewOpenConnections(
			dbcheck.WithOpenConnectionsName("test-connections-check"),
			dbcheck.WithOpenConnectionsDB(mockDB),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, "test-connections-check", check.GetName())
		assert.Equal(t, 10, result.ObservedValue)
		mockDB.AssertExpectations(t)
	})

	t.Run("check warns when connections approach maximum", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabaseStatsProvider{}
		stats := sql.DBStats{
			OpenConnections:    85,
			MaxOpenConnections: 100,
		}
		mockDB.On("Stats").Return(stats)

		check := dbcheck.NewOpenConnections(
			dbcheck.WithOpenConnectionsName("test-connections-check"),
			dbcheck.WithOpenConnectionsDB(mockDB),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusWarn, result.Status)
		assert.Contains(t, result.Output, "approaching maximum")
		assert.Equal(t, 85, result.ObservedValue)
		mockDB.AssertExpectations(t)
	})

	t.Run("check fails when connections exceed maximum", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabaseStatsProvider{}
		stats := sql.DBStats{
			OpenConnections:    100, // Equals failure threshold (100% of 100)
			MaxOpenConnections: 100,
		}
		mockDB.On("Stats").Return(stats)

		check := dbcheck.NewOpenConnections(
			dbcheck.WithOpenConnectionsName("test-connections-check"),
			dbcheck.WithOpenConnectionsDB(mockDB),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "exceed failure threshold")
		assert.Equal(t, 100, result.ObservedValue)
		mockDB.AssertExpectations(t)
	})

	t.Run("check fails when database connection is nil", func(t *testing.T) {
		t.Parallel()

		check := dbcheck.NewOpenConnections(
			dbcheck.WithOpenConnectionsName("test-connections-check"),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Equal(t, "database connection is required", result.Output)
	})

	t.Run("check works with default configuration", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabaseStatsProvider{}
		stats := sql.DBStats{
			OpenConnections:    50,
			MaxOpenConnections: 100,
		}
		mockDB.On("Stats").Return(stats)

		check := dbcheck.NewOpenConnections(dbcheck.WithOpenConnectionsDB(mockDB))
		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, "database:connections", check.GetName())
		mockDB.AssertExpectations(t)
	})

	t.Run("check with custom warn threshold", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabaseStatsProvider{}
		stats := sql.DBStats{
			OpenConnections:    60,
			MaxOpenConnections: 100,
		}
		mockDB.On("Stats").Return(stats)

		check := dbcheck.NewOpenConnections(
			dbcheck.WithOpenConnectionsDB(mockDB),
			dbcheck.WithOpenConnectionsWarnThreshold(50.0), // 50% warn threshold
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusWarn, result.Status)
		assert.Equal(t, 60, result.ObservedValue)
		mockDB.AssertExpectations(t)
	})

	t.Run("check with custom fail threshold", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabaseStatsProvider{}
		stats := sql.DBStats{
			OpenConnections:    80,
			MaxOpenConnections: 100,
		}
		mockDB.On("Stats").Return(stats)

		check := dbcheck.NewOpenConnections(
			dbcheck.WithOpenConnectionsDB(mockDB),
			dbcheck.WithOpenConnectionsFailThreshold(80.0),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "exceed failure threshold")
		mockDB.AssertExpectations(t)
	})

	t.Run("check with warn threshold", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabaseStatsProvider{}
		stats := sql.DBStats{
			OpenConnections:    70,
			MaxOpenConnections: 100,
		}
		mockDB.On("Stats").Return(stats)

		check := dbcheck.NewOpenConnections(
			dbcheck.WithOpenConnectionsDB(mockDB),
			dbcheck.WithOpenConnectionsWarnThreshold(60.0),
			dbcheck.WithOpenConnectionsFailThreshold(80.0),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusWarn, result.Status)
		assert.Equal(t, 70, result.ObservedValue)
		assert.Contains(t, result.Output, "approaching maximum")
		mockDB.AssertExpectations(t)
	})

	t.Run("check automatically infers max connections from database stats", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabaseStatsProvider{}
		stats := sql.DBStats{
			OpenConnections:    75,
			MaxOpenConnections: 200,
		}
		mockDB.On("Stats").Return(stats)

		check := dbcheck.NewOpenConnections(dbcheck.WithOpenConnectionsDB(mockDB))
		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, 75, result.ObservedValue)
		mockDB.AssertExpectations(t)
	})

	t.Run("check passes with unlimited connections", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabaseStatsProvider{}
		stats := sql.DBStats{
			OpenConnections:    50,
			MaxOpenConnections: 0, // Unlimited connections
		}
		mockDB.On("Stats").Return(stats)

		check := dbcheck.NewOpenConnections(dbcheck.WithOpenConnectionsDB(mockDB))
		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		mockDB.AssertExpectations(t)
	})
}
