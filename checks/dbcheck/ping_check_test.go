package dbcheck_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/brpaz/go-healthcheck/v2/checks"
	"github.com/brpaz/go-healthcheck/v2/checks/dbcheck"
)

// MockDatabasePinger is a mock implementation of the database pinger interface
type MockDatabasePinger struct {
	mock.Mock
}

func (m *MockDatabasePinger) PingContext(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestPingCheck_Run(t *testing.T) {
	t.Parallel()

	t.Run("check succeeds when ping succeeds", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabasePinger{}
		mockDB.On("PingContext", mock.Anything).Return(nil)

		check := dbcheck.NewPing(
			dbcheck.WithPingName("test-db-check"),
			dbcheck.WithPingDB(mockDB),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, "test-db-check", check.GetName())
		mockDB.AssertExpectations(t)
	})

	t.Run("check fails when ping fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabasePinger{}
		mockDB.On("PingContext", mock.Anything).Return(errors.New("connection error"))

		check := dbcheck.NewPing(
			dbcheck.WithPingName("test-db-check"),
			dbcheck.WithPingDB(mockDB),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "database ping failed")
		assert.Contains(t, result.Output, "connection error")
		mockDB.AssertExpectations(t)
	})

	t.Run("check fails when database connection is nil", func(t *testing.T) {
		t.Parallel()

		check := dbcheck.NewPing(
			dbcheck.WithPingName("test-db-check"),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Equal(t, "database connection is required", result.Output)
	})
}
