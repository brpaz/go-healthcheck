package pingcheck_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/brpaz/go-healthcheck/checks"
	"github.com/brpaz/go-healthcheck/checks/dbcheck/pingcheck"
)

// MockDatabasePinger is a mock implementation of the database interface
type MockDatabasePinger struct {
	mock.Mock
}

func (m *MockDatabasePinger) PingContext(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDatabasePinger) Stats() sql.DBStats {
	args := m.Called()
	return args.Get(0).(sql.DBStats)
}

func TestPingCheck_Run(t *testing.T) {
	t.Parallel()

	t.Run("check succeeds when ping succeeds", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabasePinger{}
		mockDB.On("PingContext", mock.Anything).Return(nil)

		check := pingcheck.New(
			pingcheck.WithPingName("test-db-check"),
			pingcheck.WithPingDB(mockDB),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, "test-db-check", check.GetName())
		mockDB.AssertExpectations(t)
	})

	t.Run("check fails when ping fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabasePinger{}
		mockDB.On("PingContext", mock.Anything).Return(errors.New("connection failed"))

		check := pingcheck.New(
			pingcheck.WithPingName("test-db-check"),
			pingcheck.WithPingDB(mockDB),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "connection failed")
		mockDB.AssertExpectations(t)
	})
}
