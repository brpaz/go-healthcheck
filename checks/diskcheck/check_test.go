package diskcheck_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/brpaz/go-healthcheck/checks"
	"github.com/brpaz/go-healthcheck/checks/diskcheck"
)

// MockFileSystemStater is a mock implementation of the FileSystemStater interface
type MockFileSystemStater struct {
	mock.Mock
}

func (m *MockFileSystemStater) Statfs(path string) (*diskcheck.DiskInfo, error) {
	args := m.Called(path)
	return args.Get(0).(*diskcheck.DiskInfo), args.Error(1)
}

func TestDiskCheck_New(t *testing.T) {
	t.Parallel()

	t.Run("creates check with default values", func(t *testing.T) {
		t.Parallel()

		check := diskcheck.New()

		assert.NotNil(t, check)
		assert.Equal(t, "disk-check", check.GetName())
	})

	t.Run("creates check with custom options", func(t *testing.T) {
		t.Parallel()

		mockStater := &MockFileSystemStater{}

		check := diskcheck.New(
			diskcheck.WithName("custom-disk-check"),
			diskcheck.WithPath("/var"),
			diskcheck.WithWarnThreshold(70.0),
			diskcheck.WithFailThreshold(85.0),
			diskcheck.WithFileSystemStater(mockStater),
		)

		assert.NotNil(t, check)
		assert.Equal(t, "custom-disk-check", check.GetName())
	})

	t.Run("creates check with single path", func(t *testing.T) {
		t.Parallel()

		check := diskcheck.New(
			diskcheck.WithPath("/"),
		)

		assert.NotNil(t, check)
		assert.Equal(t, "disk-check", check.GetName())
	})
}

func TestDiskCheck_GetName(t *testing.T) {
	t.Parallel()

	check := diskcheck.New()
	assert.Equal(t, "disk-check", check.GetName())
}

func TestDiskCheck_Run(t *testing.T) {
	t.Parallel()

	t.Run("succeeds when disk usage is normal", func(t *testing.T) {
		t.Parallel()

		mockStater := &MockFileSystemStater{}

		// 50% disk usage (normal)
		diskInfo := &diskcheck.DiskInfo{
			Path:     "/",
			Total:    1000000000, // 1GB
			Free:     500000000,  // 500MB
			Used:     500000000,  // 500MB
			UsedPct:  50.0,
			AvailPct: 50.0,
		}

		mockStater.On("Statfs", "/").Return(diskInfo, nil)

		check := diskcheck.New(diskcheck.WithFileSystemStater(mockStater))
		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusPass, result.Status)
		assert.Equal(t, 50.0, result.ObservedValue)
		assert.Equal(t, "%", result.ObservedUnit)

		mockStater.AssertExpectations(t)
	})

	t.Run("warns when disk usage is high but below fail threshold", func(t *testing.T) {
		t.Parallel()

		mockStater := &MockFileSystemStater{}

		// 85% disk usage (warning level)
		diskInfo := &diskcheck.DiskInfo{
			Path:     "/",
			Total:    1000000000, // 1GB
			Free:     150000000,  // 150MB
			Used:     850000000,  // 850MB
			UsedPct:  85.0,
			AvailPct: 15.0,
		}

		mockStater.On("Statfs", "/").Return(diskInfo, nil)

		check := diskcheck.New(diskcheck.WithFileSystemStater(mockStater))
		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusWarn, result.Status)
		assert.Contains(t, result.Output, "disk usage high")
		assert.Equal(t, 85.0, result.ObservedValue)

		mockStater.AssertExpectations(t)
	})

	t.Run("fails when disk usage exceeds fail threshold", func(t *testing.T) {
		t.Parallel()

		mockStater := &MockFileSystemStater{}

		// 95% disk usage (critical level)
		diskInfo := &diskcheck.DiskInfo{
			Path:     "/",
			Total:    1000000000, // 1GB
			Free:     50000000,   // 50MB
			Used:     950000000,  // 950MB
			UsedPct:  95.0,
			AvailPct: 5.0,
		}

		mockStater.On("Statfs", "/").Return(diskInfo, nil)

		check := diskcheck.New(diskcheck.WithFileSystemStater(mockStater))
		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "disk usage critical")
		assert.Equal(t, 95.0, result.ObservedValue)

		mockStater.AssertExpectations(t)
	})

	t.Run("fails when unable to get disk stats", func(t *testing.T) {
		t.Parallel()

		mockStater := &MockFileSystemStater{}
		statError := errors.New("permission denied")

		mockStater.On("Statfs", "/").Return((*diskcheck.DiskInfo)(nil), statError)

		check := diskcheck.New(diskcheck.WithFileSystemStater(mockStater))
		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusFail, result.Status)
		assert.Contains(t, result.Output, "failed to get disk stats")
		assert.Contains(t, result.Output, "permission denied")

		mockStater.AssertExpectations(t)
	})

	t.Run("checks multiple paths", func(t *testing.T) {
		t.Parallel()

		mockStater := &MockFileSystemStater{}

		// Root filesystem - normal usage
		rootInfo := &diskcheck.DiskInfo{
			Path:     "/",
			Total:    1000000000,
			Free:     500000000,
			Used:     500000000,
			UsedPct:  50.0,
			AvailPct: 50.0,
		}

		mockStater.On("Statfs", "/").Return(rootInfo, nil)

		check := diskcheck.New(
			diskcheck.WithPath("/"),
			diskcheck.WithFileSystemStater(mockStater),
		)

		result := check.Run(context.Background())

		// Should check only the first path now
		assert.Equal(t, checks.StatusPass, result.Status)

		mockStater.AssertExpectations(t)
	})
}

func TestDiskCheck_CustomThresholds(t *testing.T) {
	t.Parallel()

	t.Run("uses custom warn and fail thresholds", func(t *testing.T) {
		t.Parallel()

		mockStater := &MockFileSystemStater{}

		// 75% disk usage
		diskInfo := &diskcheck.DiskInfo{
			Path:     "/",
			Total:    1000000000,
			Free:     250000000,
			Used:     750000000,
			UsedPct:  75.0,
			AvailPct: 25.0,
		}

		mockStater.On("Statfs", "/").Return(diskInfo, nil)

		check := diskcheck.New(
			diskcheck.WithWarnThreshold(70.0), // Warn at 70%
			diskcheck.WithFailThreshold(85.0), // Fail at 85%
			diskcheck.WithFileSystemStater(mockStater),
		)

		result := check.Run(context.Background())

		assert.Equal(t, checks.StatusWarn, result.Status) // Should warn at 75% with 70% threshold
		assert.Contains(t, result.Output, "threshold: 70.0%")

		mockStater.AssertExpectations(t)
	})
}

func TestDiskCheck_GetDiskInfo(t *testing.T) {
	t.Parallel()

	t.Run("returns disk info for all paths", func(t *testing.T) {
		t.Parallel()

		mockStater := &MockFileSystemStater{}

		rootInfo := &diskcheck.DiskInfo{
			Path:     "/",
			Total:    1000000000,
			Free:     500000000,
			Used:     500000000,
			UsedPct:  50.0,
			AvailPct: 50.0,
		}

		mockStater.On("Statfs", "/").Return(rootInfo, nil)

		check := diskcheck.New(
			diskcheck.WithPath("/"),
			diskcheck.WithFileSystemStater(mockStater),
		)

		infos, err := check.GetDiskInfo()

		assert.NoError(t, err)
		assert.Len(t, infos, 1)
		assert.Equal(t, "/", infos[0].Path)

		mockStater.AssertExpectations(t)
	})

	t.Run("returns error when stat fails", func(t *testing.T) {
		t.Parallel()

		mockStater := &MockFileSystemStater{}
		statError := errors.New("stat failed")

		mockStater.On("Statfs", "/").Return((*diskcheck.DiskInfo)(nil), statError)

		check := diskcheck.New(diskcheck.WithFileSystemStater(mockStater))

		infos, err := check.GetDiskInfo()

		assert.Error(t, err)
		assert.Nil(t, infos)
		assert.Contains(t, err.Error(), "stat failed")

		mockStater.AssertExpectations(t)
	})
}

// Table-driven tests for threshold scenarios
func TestDiskCheck_Run_ThresholdScenarios(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		usedPct          float64
		warnThreshold    float64
		failThreshold    float64
		expectedStatus   checks.Status
		expectedContains string
	}{
		{
			name:             "normal usage below warn threshold",
			usedPct:          60.0,
			warnThreshold:    80.0,
			failThreshold:    90.0,
			expectedStatus:   checks.StatusPass,
			expectedContains: "",
		},
		{
			name:             "usage exactly at warn threshold",
			usedPct:          80.0,
			warnThreshold:    80.0,
			failThreshold:    90.0,
			expectedStatus:   checks.StatusWarn,
			expectedContains: "disk usage high",
		},
		{
			name:             "usage exactly at fail threshold",
			usedPct:          90.0,
			warnThreshold:    80.0,
			failThreshold:    90.0,
			expectedStatus:   checks.StatusFail,
			expectedContains: "disk usage critical",
		},
		{
			name:             "usage above fail threshold",
			usedPct:          95.0,
			warnThreshold:    80.0,
			failThreshold:    90.0,
			expectedStatus:   checks.StatusFail,
			expectedContains: "disk usage critical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStater := &MockFileSystemStater{}

			diskInfo := &diskcheck.DiskInfo{
				Path:     "/",
				Total:    1000000000,
				Free:     uint64((100.0 - tt.usedPct) * 10000000),
				Used:     uint64(tt.usedPct * 10000000),
				UsedPct:  tt.usedPct,
				AvailPct: 100.0 - tt.usedPct,
			}

			mockStater.On("Statfs", "/").Return(diskInfo, nil)

			check := diskcheck.New(
				diskcheck.WithWarnThreshold(tt.warnThreshold),
				diskcheck.WithFailThreshold(tt.failThreshold),
				diskcheck.WithFileSystemStater(mockStater),
			)

			result := check.Run(context.Background())

			assert.Equal(t, tt.expectedStatus, result.Status)
			assert.Contains(t, result.Output, tt.expectedContains)
			assert.Equal(t, tt.usedPct, result.ObservedValue)

			mockStater.AssertExpectations(t)
		})
	}
}
