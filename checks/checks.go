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
	ComponentID   string    `json:"component_id,omitempty"`
	ComponentType string    `json:"component_type,omitempty"`
}

// Check is an interface that any health check implementation must satisfy.
type Check interface {
	GetName() string
	Run(ctx context.Context) []Result
}
