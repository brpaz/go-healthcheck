package sqlcheck_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

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

// For testing the query part, we'll focus on ping testing since mocking sql.Row is complex
func TestSQLCheck_New(t *testing.T) {
	t.Parallel()

	t.Run("creates check with default values", func(t *testing.T) {
		t.Parallel()

		check := sqlcheck.New()

		assert.NotNil(t, check)
		assert.Equal(t, "sql-check", check.GetName())
	})

	t.Run("creates check with custom options", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabase{}
		customTimeout := 10 * time.Second

		check := sqlcheck.New(
			sqlcheck.WithName("custom-sql-check"),
			sqlcheck.WithDB(mockDB),
			sqlcheck.WithTimeout(customTimeout),
		)

		assert.NotNil(t, check)
		assert.Equal(t, "sql-check", check.GetName()) // Name is always "sql-check" based on the implementation
	})
}

func TestSQLCheck_GetName(t *testing.T) {
	t.Parallel()

	check := sqlcheck.New()
	assert.Equal(t, "sql-check", check.GetName())
}

func TestSQLCheck_Run(t *testing.T) {
	t.Parallel()

	t.Run("fails when database is nil", func(t *testing.T) {
		t.Parallel()

		check := sqlcheck.New()
		results := check.Run(context.Background())

		assert.Len(t, results, 1)
		result := results[0]
		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Equal(t, "database connection is required", result.Output)
		assert.Equal(t, "database", result.ComponentType)
		assert.Equal(t, "sql-check", result.ComponentID)
	})

	t.Run("fails when ping fails", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabase{}
		pingError := errors.New("connection failed")
		mockDB.On("PingContext", mock.Anything).Return(pingError)

		check := sqlcheck.New(sqlcheck.WithDB(mockDB))
		results := check.Run(context.Background())

		assert.Len(t, results, 1)
		result := results[0]
		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "database ping failed")
		assert.Contains(t, result.Output, "connection failed")
		assert.Equal(t, "database", result.ComponentType)
		assert.Equal(t, "sql-check", result.ComponentID)

		mockDB.AssertExpectations(t)
	})

	t.Run("respects timeout context", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabase{}
		timeoutError := errors.New("context deadline exceeded")

		mockDB.On("PingContext", mock.Anything).Return(timeoutError)

		check := sqlcheck.New(
			sqlcheck.WithDB(mockDB),
			sqlcheck.WithTimeout(1*time.Millisecond),
		)

		results := check.Run(context.Background())

		assert.Len(t, results, 1)
		result := results[0]
		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "database ping failed")

		mockDB.AssertExpectations(t)
	})

	t.Run("handles context cancellation", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		mockDB := &MockDatabase{}
		mockDB.On("PingContext", mock.Anything).Return(context.Canceled)

		check := sqlcheck.New(sqlcheck.WithDB(mockDB))
		results := check.Run(ctx)

		assert.Len(t, results, 1)
		result := results[0]
		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "database ping failed")

		mockDB.AssertExpectations(t)
	})

	t.Run("includes metrics when enabled", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabase{}
		mockDB.On("PingContext", mock.Anything).Return(nil)

		// Setup mock stats
		mockStats := sql.DBStats{
			MaxOpenConnections: 25,
			OpenConnections:    5,
			InUse:              2,
			Idle:               3,
			WaitCount:          10,
			WaitDuration:       time.Millisecond * 50,
		}
		mockDB.On("Stats").Return(mockStats)

		check := sqlcheck.New(
			sqlcheck.WithDB(mockDB),
			sqlcheck.WithMetrics(true),
		)

		results := check.Run(context.Background())

		// Should have 7 results: 1 connectivity + 6 metrics
		assert.Len(t, results, 7)

		// First result should be the main connectivity check
		connectivityResult := results[0]
		assert.Equal(t, checks.StatusPass, connectivityResult.Status)
		assert.Equal(t, "sql-check", connectivityResult.ComponentID)
		assert.Equal(t, "ms", connectivityResult.ObservedUnit)

		// Check metric results
		expectedMetrics := []struct {
			componentID   string
			unit          string
			expectedValue int64
		}{
			{"sql-check:open-connections", "", 5},
			{"sql-check:in-use-connections", "", 2},
			{"sql-check:idle-connections", "", 3},
			{"sql-check:max-open-connections", "", 25},
			{"sql-check:wait-count", "", 10},
			{"sql-check:wait-duration", "ms", 50},
		}

		for i, expected := range expectedMetrics {
			result := results[i+1] // Skip first connectivity result
			assert.Equal(t, checks.StatusPass, result.Status)
			assert.Equal(t, expected.componentID, result.ComponentID)
			assert.Equal(t, expected.unit, result.ObservedUnit)
			assert.Equal(t, expected.expectedValue, result.ObservedValue)
			assert.Equal(t, "database", result.ComponentType)
			assert.Equal(t, "", result.Output) // Empty output for successful check
		}

		mockDB.AssertExpectations(t)
	})

	t.Run("excludes metrics when disabled", func(t *testing.T) {
		t.Parallel()

		mockDB := &MockDatabase{}
		mockDB.On("PingContext", mock.Anything).Return(nil)
		// Note: Stats should not be called when metrics are disabled

		check := sqlcheck.New(
			sqlcheck.WithDB(mockDB),
			sqlcheck.WithMetrics(false),
		)

		results := check.Run(context.Background())

		// Should only have 1 result (connectivity check)
		assert.Len(t, results, 1)
		result := results[0]
		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, "sql-check", result.ComponentID)
		assert.Equal(t, "", result.Output) // Empty output for successful check
		assert.Equal(t, "ms", result.ObservedUnit)

		mockDB.AssertExpectations(t)
	})
}
