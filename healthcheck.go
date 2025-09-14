// Package healthcheck provides a library for creating health check endpoints.
// It includes various built-in checks and new ones can be added easily.
// It tooks inspiration from https://inadarei.github.io/rfc-healthcheck/ for the response structure.
// Example Usage:
//
// package main
//
// import "net/http"
// import "github.com/brpaz/go-healthcheck"
// import "github.com/brpaz/go-healthcheck/pkg/checks/mockcheck"
//
//	func main() {
//	    mycheck := mockcheck.New(
//			mockcheck.WithName("my-check"),
//		)
//		hc := healthcheck.NewHealthChecker(
//			healthcheck.WithServiceID("my-service"),
//			healthcheck.WithDescription("My Service"),
//			healthcheck.WithVersion("1.0.0"),
//			healthcheck.WithReleaseID("1.0.0-SNAPSHOT"),
//			healthcheck.WithCheck(mycheck),
//		)
//
//		http.HandleFunc("/health", healthcheck.Handler(hc))
//		http.ListenAndServe(":8080", nil)
//	}
package healthcheck

import (
	"context"

	"github.com/brpaz/go-healthcheck/checks"
)

type HealthChecker interface {
	Execute(ctx context.Context) CheckRunResult
}

// HealthCheck aggregates multiple healthchecks and provides metadata about the service.
type HealthCheck struct {
	ServiceID   string
	Description string
	Version     string
	ReleaseID   string
	Checks      []checks.Check
}

// Option is a functional option for configuring HealthCheck.
type Option func(*HealthCheck)

// WithServiceID sets the service ID.
func WithServiceID(id string) Option {
	return func(h *HealthCheck) {
		h.ServiceID = id
	}
}

// WithDescription sets the description.
func WithDescription(desc string) Option {
	return func(h *HealthCheck) {
		h.Description = desc
	}
}

// WithVersion sets the version.
func WithVersion(version string) Option {
	return func(h *HealthCheck) {
		h.Version = version
	}
}

// WithReleaseID sets the release ID.
func WithReleaseID(id string) Option {
	return func(h *HealthCheck) {
		h.ReleaseID = id
	}
}

// WithCheck registers a check in the HealthCheck.
func WithCheck(check checks.Check) Option {
	return func(h *HealthCheck) {
		h.Checks = append(h.Checks, check)
	}
}

// NewHealthCheck creates a new HealthChecker instance the provided options.
func NewHealthCheck(opts ...Option) *HealthCheck {
	h := &HealthCheck{
		Checks: make([]checks.Check, 0),
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// AddCheck adds a new check to the HealthCheck instance.
func (h *HealthCheck) AddCheck(check checks.Check) {
	h.Checks = append(h.Checks, check)
}

// GetChecks returns the registered checks.
func (h *HealthCheck) GetChecks() []checks.Check {
	return h.Checks
}

// CheckRunResult aggregates the result of running a group of checks.
type CheckRunResult struct {
	Status checks.Status
	Checks map[string][]checks.Result
}

// Execute runs all registered healthchecks and returns an aggregated result, composed of the
// overall status and the individual results of each check.
// The final status is determined as follows:
// - If any check returns StatusFail, the overall status is StatusFail.
// - If no checks return StatusFail but at least one returns StatusWarn, the overall status is StatusWarn.
// - If all checks return StatusPass, the overall status is StatusPass.
func (h *HealthCheck) Execute(ctx context.Context) CheckRunResult {
	type resultCollector struct {
		name   string
		result []checks.Result
	}

	resultsChan := make(chan resultCollector, len(h.Checks))

	for _, check := range h.Checks {
		go func(c checks.Check) {
			result := c.Run(ctx)
			resultsChan <- resultCollector{
				name:   c.GetName(),
				result: result,
			}
		}(check)
	}

	// Collect results
	results := make(map[string][]checks.Result)
	status := checks.StatusPass

	for range h.Checks {
		cr := <-resultsChan
		results[cr.name] = append(results[cr.name], cr.result...)

		for _, result := range cr.result {
			if result.Status == checks.StatusFail {
				status = checks.StatusFail
			} else if result.Status == checks.StatusWarn && status != checks.StatusFail {
				status = checks.StatusWarn
			}
		}
	}

	return CheckRunResult{
		Status: status,
		Checks: results,
	}
}
