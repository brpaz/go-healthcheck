package sqlcheck_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/brpaz/go-healthcheck/checks"
	"github.com/brpaz/go-healthcheck/checks/sqlcheck"
)

// MockDatabase is a mock implementation of the database interface
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) PingContext(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDatabase) Stats() sql.DBStats {
	args := m.Called()
	return args.Get(0).(sql.DBStats)
}

func TestConnectivityCheck_Run(t *testing.T) {
	t.Parallel()

	t.Run("connectivity check succeeds when ping succeeds", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabase{}
		mockDB.On("PingContext", mock.Anything).Return(nil)

		check := sqlcheck.NewConnectivityCheck(
			sqlcheck.WithConnectivityName("test-db-check"),
			sqlcheck.WithConnectivityDB(mockDB),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, "test-db-check", check.GetName())
		mockDB.AssertExpectations(t)
	})

	t.Run("connectivity check fails when ping fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabase{}
		mockDB.On("PingContext", mock.Anything).Return(errors.New("connection failed"))

		check := sqlcheck.NewConnectivityCheck(
			sqlcheck.WithConnectivityName("test-db-check"),
			sqlcheck.WithConnectivityDB(mockDB),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "connection failed")
		mockDB.AssertExpectations(t)
	})
}

func TestOpenConnectionsCheck_Run(t *testing.T) {
	t.Parallel()

	t.Run("open connections metric check succeeds", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabase{}
		stats := sql.DBStats{
			OpenConnections: 5,
		}
		mockDB.On("Stats").Return(stats)

		check := sqlcheck.NewOpenConnectionsCheck("test-db-check", mockDB)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, "test-db-check:open-connections", check.GetName())
		assert.Equal(t, int64(5), result.ObservedValue)
		mockDB.AssertExpectations(t)
	})
}

func TestInUseConnectionsCheck_Run(t *testing.T) {
	t.Parallel()

	t.Run("in-use connections metric check succeeds", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabase{}
		stats := sql.DBStats{
			InUse: 3,
		}
		mockDB.On("Stats").Return(stats)

		check := sqlcheck.NewInUseConnectionsCheck("test-db-check", mockDB)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, "test-db-check:in-use-connections", check.GetName())
		assert.Equal(t, int64(3), result.ObservedValue)
		mockDB.AssertExpectations(t)
	})
}
