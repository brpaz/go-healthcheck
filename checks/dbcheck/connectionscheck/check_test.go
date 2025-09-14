package connectionscheck_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/brpaz/go-healthcheck/v2/checks"
	"github.com/brpaz/go-healthcheck/v2/checks/dbcheck/connectionscheck"
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

		check := connectionscheck.New(
			connectionscheck.WithName("test-connections-check"),
			connectionscheck.WithDB(mockDB),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, "test-connections-check", check.GetName())
		assert.Equal(t, 10, result.ObservedValue)
		assert.Contains(t, result.Output, "open connections: 10/100")
		mockDB.AssertExpectations(t)
	})

	t.Run("check warns when connections approach maximum", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabaseStatsProvider{}
		stats := sql.DBStats{
			OpenConnections:    80, // 80% of 100 (equals default warn threshold)
			MaxOpenConnections: 100,
		}
		mockDB.On("Stats").Return(stats)

		check := connectionscheck.New(
			connectionscheck.WithName("test-connections-check"),
			connectionscheck.WithDB(mockDB),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusWarn, result.Status)
		assert.Equal(t, 80, result.ObservedValue)
		assert.Contains(t, result.Output, "approaching maximum")
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

		check := connectionscheck.New(
			connectionscheck.WithName("test-connections-check"),
			connectionscheck.WithDB(mockDB),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Equal(t, 100, result.ObservedValue)
		assert.Contains(t, result.Output, "exceed failure threshold")
		mockDB.AssertExpectations(t)
	})

	t.Run("check fails when database connection is nil", func(t *testing.T) {
		t.Parallel()

		check := connectionscheck.New(
			connectionscheck.WithName("test-connections-check"),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "database connection is required")
	})

	t.Run("check works with default configuration", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabaseStatsProvider{}
		stats := sql.DBStats{
			OpenConnections:    50,
			MaxOpenConnections: 100,
		}
		mockDB.On("Stats").Return(stats)

		check := connectionscheck.New(
			connectionscheck.WithDB(mockDB),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, "db-connections-check", check.GetName())
		assert.Equal(t, 50, result.ObservedValue)
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

		check := connectionscheck.New(
			connectionscheck.WithDB(mockDB),
			connectionscheck.WithWarnThreshold(0.5), // 50% warn threshold
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusWarn, result.Status)
		assert.Equal(t, 60, result.ObservedValue)
		assert.Contains(t, result.Output, "approaching maximum")
		mockDB.AssertExpectations(t)
	})

	t.Run("check with custom fail threshold", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabaseStatsProvider{}
		stats := sql.DBStats{
			OpenConnections:    85,
			MaxOpenConnections: 100,
		}
		mockDB.On("Stats").Return(stats)

		check := connectionscheck.New(
			connectionscheck.WithDB(mockDB),
			connectionscheck.WithFailThreshold(0.8),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Equal(t, 85, result.ObservedValue)
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

		check := connectionscheck.New(
			connectionscheck.WithDB(mockDB),
			connectionscheck.WithWarnThreshold(0.6),
			connectionscheck.WithFailThreshold(0.8),
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
			OpenConnections:    60,
			MaxOpenConnections: 100,
		}
		mockDB.On("Stats").Return(stats)

		check := connectionscheck.New(
			connectionscheck.WithDB(mockDB),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, 60, result.ObservedValue)
		assert.Contains(t, result.Output, "open connections: 60/100")
		mockDB.AssertExpectations(t)
	})

	t.Run("check passes with unlimited connections", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabaseStatsProvider{}
		stats := sql.DBStats{
			OpenConnections:    10,
			MaxOpenConnections: 0, // Unlimited connections
		}
		mockDB.On("Stats").Return(stats)

		check := connectionscheck.New(
			connectionscheck.WithDB(mockDB),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		mockDB.AssertExpectations(t)
	})
}
