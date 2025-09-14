package diskcheck

import (
	"fmt"
	"syscall"
)

// FileSystemStater defines the interface for getting filesystem statistics
type FileSystemStater interface {
	Statfs(path string) (*DiskInfo, error)
}

// DefaultFileSystemStater implements FileSystemStater using syscall.Statfs
type DefaultFileSystemStater struct{}

func (d *DefaultFileSystemStater) Statfs(path string) (*DiskInfo, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return nil, fmt.Errorf("failed to get filesystem stats for %s: %w", path, err)
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	used := total - (stat.Bfree * uint64(stat.Bsize))

	var usedPct, availPct float64
	if total > 0 {
		usedPct = float64(used) / float64(total) * 100
		availPct = float64(free) / float64(total) * 100
	}

	return &DiskInfo{
		Path:     path,
		Total:    total,
		Free:     free,
		Used:     used,
		UsedPct:  usedPct,
		AvailPct: availPct,
	}, nil
}
