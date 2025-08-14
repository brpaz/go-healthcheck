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
//		hc := healthcheck.New(
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
	"encoding/json"
	"net/http"

	"github.com/brpaz/go-healthcheck/pkg/checks"
)

// HealthCheck represents a collection of health checks for a service.
type HealthCheck struct {
	ServiceID   string
	Description string
	Version     string
	ReleaseID   string
	Checks      []checks.Check
}

// Response represents the health check response structure.
// Health Check Response Format for HTTP APIs uses JSON format described in RFC 8259 and has the media type "application/health+json".
// Its content consists of a single mandatory root field ("status") and several optional fields:
// See https://tools.ietf.org/id/draft-inadarei-api-health-check-05.html#section-3
type Response struct {
	ServiceID   string                     `json:"service_id,omitempty"`
	Description string                     `json:"description,omitempty"`
	Version     string                     `json:"version,omitempty"`
	ReleaseID   string                     `json:"release_id,omitempty"`
	Output      string                     `json:"output,omitempty"`
	Status      checks.Status              `json:"status"`
	Checks      map[string][]checks.Result `json:"checks"`
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

// New creates a new HealthCheck instance with optional configuration.
func New(opts ...Option) *HealthCheck {
	h := &HealthCheck{
		Checks: make([]checks.Check, 0),
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// Execute runs all checks in parallel and returns the result.
// The global health check status is calculated based on the statuses of all checks
func (h *HealthCheck) Execute(ctx context.Context) Response {
	type checkResult struct {
		name   string
		result []checks.Result
	}

	resultsChan := make(chan checkResult, len(h.Checks))

	for _, check := range h.Checks {
		go func(c checks.Check) {
			result := c.Run(ctx)
			resultsChan <- checkResult{
				name:   c.GetName(),
				result: result,
			}
		}(check)
	}

	// Collect results
	checksResults := make(map[string][]checks.Result)
	status := checks.StatusPass

	for range h.Checks {
		cr := <-resultsChan
		checksResults[cr.name] = append(checksResults[cr.name], cr.result...)

		for _, result := range cr.result {
			if result.Status == checks.StatusFail {
				status = checks.StatusFail
			} else if result.Status == checks.StatusWarn && status != checks.StatusFail {
				status = checks.StatusWarn
			}
		}
	}

	response := Response{
		ServiceID:   h.ServiceID,
		Description: h.Description,
		Version:     h.Version,
		ReleaseID:   h.ReleaseID,
		Status:      status,
		Checks:      checksResults,
	}

	return response
}

// Handler provides an HTTP handler that can be used to serve the health check endpoint.
func Handler(hc *HealthCheck) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		response := hc.Execute(ctx)
		w.Header().Set("Content-Type", "application/health+json")

		if response.Status == checks.StatusFail {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		_ = json.NewEncoder(w).Encode(response)
	}
}
