package checks

import (
	"context"
	"time"
)

// Status represents the status of a health check.
type Status string

const (
	StatusPass Status = "pass"
	StatusFail Status = "fail"
	StatusWarn Status = "warn"
)

// Result represents the result of an individual health check execution.
type Result struct {
	Status        Status    `json:"status"`
	Output        string    `json:"output,omitempty"`
	Time          time.Time `json:"time"`
	ObservedValue any       `json:"observed_value,omitempty"`
	ObservedUnit  string    `json:"observed_unit,omitempty"`
}

// Check is an interface that any health check implementation must satisfy.
// Each Check represents a single measurement/test and should return exactly one Result.
// For components that need multiple measurements (e.g., database connection + metrics),
// create separate Check implementations for each measurement.
type Check interface {
	GetName() string                // Returns the unique name for this specific check (e.g., "db-check:open-connections")
	Run(ctx context.Context) Result // Returns a single result for this specific check
}
