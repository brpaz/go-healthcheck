package tcpcheck_test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/brpaz/go-healthcheck/v2/checks"
	"github.com/brpaz/go-healthcheck/v2/checks/tcpcheck"
)

// MockDialer is a mock implementation of the Dialer interface
type MockDialer struct {
	mock.Mock
}

func (m *MockDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	args := m.Called(ctx, network, address)
	return args.Get(0).(net.Conn), args.Error(1)
}

// MockConn is a mock implementation of net.Conn
type MockConn struct {
	mock.Mock
}

func (m *MockConn) Read(b []byte) (n int, err error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *MockConn) Write(b []byte) (n int, err error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *MockConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockConn) LocalAddr() net.Addr {
	args := m.Called()
	return args.Get(0).(net.Addr)
}

func (m *MockConn) RemoteAddr() net.Addr {
	args := m.Called()
	return args.Get(0).(net.Addr)
}

func (m *MockConn) SetDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockConn) SetReadDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockConn) SetWriteDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func TestTCPCheck_New(t *testing.T) {
	t.Parallel()

	t.Run("creates check with default values", func(t *testing.T) {
		t.Parallel()

		check := tcpcheck.New()

		assert.NotNil(t, check)
		assert.Equal(t, "tcp-check", check.GetName())
	})

	t.Run("creates check with custom options", func(t *testing.T) {
		t.Parallel()

		mockDialer := &MockDialer{}
		customTimeout := 10 * time.Second

		check := tcpcheck.New(
			tcpcheck.WithName("custom-tcp-check"),
			tcpcheck.WithHost("localhost"),
			tcpcheck.WithPort(8080),
			tcpcheck.WithNetwork(tcpcheck.TCP),
			tcpcheck.WithTimeout(customTimeout),
			tcpcheck.WithDialer(mockDialer),
		)

		assert.NotNil(t, check)
		assert.Equal(t, "custom-tcp-check", check.GetName())
		assert.Equal(t, "tcp://localhost:8080", check.Address())
	})
}

func TestTCPCheck_GetName(t *testing.T) {
	t.Parallel()

	check := tcpcheck.New()
	assert.Equal(t, "tcp-check", check.GetName())
}

func TestTCPCheck_Address(t *testing.T) {
	t.Parallel()

	t.Run("returns correct TCP address", func(t *testing.T) {
		t.Parallel()

		check := tcpcheck.New(
			tcpcheck.WithHost("example.com"),
			tcpcheck.WithPort(443),
			tcpcheck.WithNetwork(tcpcheck.TCP),
		)

		assert.Equal(t, "tcp://example.com:443", check.Address())
	})

	t.Run("returns correct UDP address", func(t *testing.T) {
		t.Parallel()

		check := tcpcheck.New(
			tcpcheck.WithHost("localhost"),
			tcpcheck.WithPort(53),
			tcpcheck.WithNetwork(tcpcheck.UDP),
		)

		assert.Equal(t, "udp://localhost:53", check.Address())
	})
}

func TestTCPCheck_Run(t *testing.T) {
	t.Parallel()

	t.Run("fails when host is empty", func(t *testing.T) {
		t.Parallel()

		check := tcpcheck.New(tcpcheck.WithPort(8080))
		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Equal(t, "host is required", result.Output)
	})

	t.Run("fails when port is invalid", func(t *testing.T) {
		t.Parallel()

		check := tcpcheck.New(
			tcpcheck.WithHost("localhost"),
			tcpcheck.WithPort(0),
		)
		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "invalid port")
	})

	t.Run("fails when port is too high", func(t *testing.T) {
		t.Parallel()

		check := tcpcheck.New(
			tcpcheck.WithHost("localhost"),
			tcpcheck.WithPort(70000),
		)
		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "invalid port")
	})

	t.Run("succeeds when connection is successful", func(t *testing.T) {
		t.Parallel()

		mockDialer := &MockDialer{}
		mockConn := &MockConn{}

		mockDialer.On("DialContext", mock.Anything, "tcp", "localhost:8080").Return(mockConn, nil)
		mockConn.On("Close").Return(nil)

		check := tcpcheck.New(
			tcpcheck.WithHost("localhost"),
			tcpcheck.WithPort(8080),
			tcpcheck.WithDialer(mockDialer),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, "ms", result.ObservedUnit)
		assert.GreaterOrEqual(t, result.ObservedValue, int64(0))

		mockDialer.AssertExpectations(t)
		mockConn.AssertExpectations(t)
	})

	t.Run("fails when connection fails", func(t *testing.T) {
		t.Parallel()

		mockDialer := &MockDialer{}
		connError := errors.New("connection refused")

		mockDialer.On("DialContext", mock.Anything, "tcp", "localhost:8080").Return((*MockConn)(nil), connError)

		check := tcpcheck.New(
			tcpcheck.WithHost("localhost"),
			tcpcheck.WithPort(8080),
			tcpcheck.WithDialer(mockDialer),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "failed to connect")
		assert.Contains(t, result.Output, "connection refused")

		mockDialer.AssertExpectations(t)
	})

	t.Run("succeeds even if close fails", func(t *testing.T) {
		t.Parallel()

		mockDialer := &MockDialer{}
		mockConn := &MockConn{}

		mockDialer.On("DialContext", mock.Anything, "tcp", "localhost:8080").Return(mockConn, nil)
		mockConn.On("Close").Return(errors.New("close failed"))

		check := tcpcheck.New(
			tcpcheck.WithHost("localhost"),
			tcpcheck.WithPort(8080),
			tcpcheck.WithDialer(mockDialer),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Contains(t, result.Output, "connection successful but failed to close")

		mockDialer.AssertExpectations(t)
		mockConn.AssertExpectations(t)
	})

	t.Run("respects timeout context", func(t *testing.T) {
		t.Parallel()

		mockDialer := &MockDialer{}
		timeoutError := errors.New("context deadline exceeded")

		mockDialer.On("DialContext", mock.Anything, "tcp", "localhost:8080").Return((*MockConn)(nil), timeoutError)

		check := tcpcheck.New(
			tcpcheck.WithHost("localhost"),
			tcpcheck.WithPort(8080),
			tcpcheck.WithTimeout(1*time.Millisecond),
			tcpcheck.WithDialer(mockDialer),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "failed to connect")

		mockDialer.AssertExpectations(t)
	})

	t.Run("handles context cancellation", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		mockDialer := &MockDialer{}
		mockDialer.On("DialContext", mock.Anything, "tcp", "localhost:8080").Return((*MockConn)(nil), context.Canceled)

		check := tcpcheck.New(
			tcpcheck.WithHost("localhost"),
			tcpcheck.WithPort(8080),
			tcpcheck.WithDialer(mockDialer),
		)

		result := check.Run(ctx)

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "failed to connect")

		mockDialer.AssertExpectations(t)
	})
}

func TestTCPCheck_Options(t *testing.T) {
	t.Parallel()

	t.Run("WithName option", func(t *testing.T) {
		t.Parallel()

		check := tcpcheck.New(tcpcheck.WithName("custom-name"))
		assert.Equal(t, "custom-name", check.GetName())
	})

	t.Run("WithHost and WithPort options", func(t *testing.T) {
		t.Parallel()

		check := tcpcheck.New(
			tcpcheck.WithHost("example.com"),
			tcpcheck.WithPort(443),
		)

		assert.Equal(t, "tcp://example.com:443", check.Address())
	})

	t.Run("WithNetwork option", func(t *testing.T) {
		t.Parallel()

		check := tcpcheck.New(
			tcpcheck.WithHost("localhost"),
			tcpcheck.WithPort(53),
			tcpcheck.WithNetwork(tcpcheck.UDP),
		)

		assert.Equal(t, "udp://localhost:53", check.Address())
	})

	t.Run("WithTimeout option", func(t *testing.T) {
		t.Parallel()

		mockDialer := &MockDialer{}
		customTimeout := 100 * time.Millisecond

		mockDialer.On("DialContext", mock.Anything, "tcp", "localhost:8080").Run(func(args mock.Arguments) {
			ctx := args.Get(0).(context.Context)
			deadline, ok := ctx.Deadline()
			assert.True(t, ok, "Context should have a deadline")

			// Verify the timeout is approximately what we set
			expectedDeadline := time.Now().Add(customTimeout)
			assert.WithinDuration(t, expectedDeadline, deadline, 50*time.Millisecond)
		}).Return((*MockConn)(nil), errors.New("timeout"))

		check := tcpcheck.New(
			tcpcheck.WithHost("localhost"),
			tcpcheck.WithPort(8080),
			tcpcheck.WithTimeout(customTimeout),
			tcpcheck.WithDialer(mockDialer),
		)

		result := check.Run(context.Background())
		assert.Equal(t, checks.StatusFail, result.Status)

		mockDialer.AssertExpectations(t)
	})
}
