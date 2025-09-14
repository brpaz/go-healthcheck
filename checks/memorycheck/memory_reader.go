package memorycheck

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MemoryReader interface for reading memory stats (useful for testing)
type MemoryReader interface {
	ReadMemoryStats() (*MemoryStats, error)
}

// DefaultMemoryReader reads memory stats from /proc/meminfo
type DefaultMemoryReader struct{}

// ReadMemoryStats reads memory statistics from /proc/meminfo
func (r *DefaultMemoryReader) ReadMemoryStats() (*MemoryStats, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, fmt.Errorf("failed to open /proc/meminfo: %w", err)
	}
	defer func() { _ = file.Close() }()

	var memTotal, memAvailable uint64
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			if val, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
				memTotal = val * 1024 // Convert from KB to bytes
			}
		case "MemAvailable:":
			if val, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
				memAvailable = val * 1024 // Convert from KB to bytes
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading /proc/meminfo: %w", err)
	}

	if memTotal == 0 {
		return nil, fmt.Errorf("could not read MemTotal from /proc/meminfo")
	}

	used := memTotal - memAvailable
	usedPct := float64(used) / float64(memTotal) * 100

	return &MemoryStats{
		Total:     memTotal,
		Available: memAvailable,
		Used:      used,
		UsedPct:   usedPct,
	}, nil
}
