package redischeck_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/brpaz/go-healthcheck/checks"
	"github.com/brpaz/go-healthcheck/checks/redischeck"
)

// MockRedisClient is a mock implementation of the RedisClient interface
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestRedisCheck_New(t *testing.T) {
	t.Parallel()

	t.Run("creates check with default values", func(t *testing.T) {
		t.Parallel()

		check := redischeck.New()

		assert.NotNil(t, check)
		assert.Equal(t, "redis-check", check.GetName())
	})

	t.Run("creates check with custom options", func(t *testing.T) {
		t.Parallel()

		mockClient := &MockRedisClient{}
		customTimeout := 10 * time.Second

		check := redischeck.New(
			redischeck.WithName("custom-redis-check"),
			redischeck.WithClient(mockClient),
			redischeck.WithTimeout(customTimeout),
			redischeck.WithComponentType("cache"),
			redischeck.WithComponentID("redis-primary"),
		)

		assert.NotNil(t, check)
		assert.Equal(t, "custom-redis-check", check.GetName())
	})
}

func TestRedisCheck_GetName(t *testing.T) {
	t.Parallel()

	check := redischeck.New()
	assert.Equal(t, "redis-check", check.GetName())
}

func TestRedisCheck_Run(t *testing.T) {
	t.Parallel()

	t.Run("fails when client is nil", func(t *testing.T) {
		t.Parallel()

		check := redischeck.New()
		results := check.Run(context.Background())

		assert.Len(t, results, 1)
		result := results[0]
		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Equal(t, "Redis client is required", result.Output)
		assert.Equal(t, "datastore", result.ComponentType)
		assert.Equal(t, "redis", result.ComponentID)
	})

	t.Run("succeeds when ping succeeds", func(t *testing.T) {
		t.Parallel()

		mockClient := &MockRedisClient{}
		mockClient.On("Ping", mock.Anything).Return(nil)

		check := redischeck.New(redischeck.WithClient(mockClient))
		results := check.Run(context.Background())

		assert.Len(t, results, 1)
		result := results[0]
		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, "datastore", result.ComponentType)
		assert.Equal(t, "redis", result.ComponentID)
		assert.Equal(t, "ms", result.ObservedUnit)
		assert.GreaterOrEqual(t, result.ObservedValue, int64(0))

		mockClient.AssertExpectations(t)
	})

	t.Run("fails when ping fails", func(t *testing.T) {
		t.Parallel()

		mockClient := &MockRedisClient{}
		pingError := errors.New("connection refused")
		mockClient.On("Ping", mock.Anything).Return(pingError)

		check := redischeck.New(redischeck.WithClient(mockClient))
		results := check.Run(context.Background())

		assert.Len(t, results, 1)
		result := results[0]
		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "Redis ping failed")
		assert.Contains(t, result.Output, "connection refused")

		mockClient.AssertExpectations(t)
	})

	t.Run("respects timeout context", func(t *testing.T) {
		t.Parallel()

		mockClient := &MockRedisClient{}
		timeoutError := errors.New("context deadline exceeded")
		mockClient.On("Ping", mock.Anything).Return(timeoutError)

		check := redischeck.New(
			redischeck.WithClient(mockClient),
			redischeck.WithTimeout(1*time.Millisecond),
		)

		results := check.Run(context.Background())

		assert.Len(t, results, 1)
		result := results[0]
		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "Redis ping failed")

		mockClient.AssertExpectations(t)
	})

	t.Run("handles context cancellation", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		mockClient := &MockRedisClient{}
		mockClient.On("Ping", mock.Anything).Return(context.Canceled)

		check := redischeck.New(redischeck.WithClient(mockClient))
		results := check.Run(ctx)

		assert.Len(t, results, 1)
		result := results[0]
		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "Redis ping failed")

		mockClient.AssertExpectations(t)
	})
}
