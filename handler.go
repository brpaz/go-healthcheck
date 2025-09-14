package healthcheck

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/brpaz/go-healthcheck/checks"
)

// HealthHttpResponse represents the structure of the health check HTTP response.
type HealthHttpResponse struct {
	ServiceID   string                     `json:"serviceId,omitempty"`
	Description string                     `json:"description,omitempty"`
	Version     string                     `json:"version,omitempty"`
	ReleaseID   string                     `json:"releaseId,omitempty"`
	Output      string                     `json:"output,omitempty"`
	Status      checks.Status              `json:"status"`
	Checks      map[string][]checks.Result `json:"checks"`
}

func buildOutput(checks map[string][]checks.Result) string {
	var outputs []string
	for checkName, results := range checks {
		for _, result := range results {
			if result.Output != "" {
				outputs = append(outputs, checkName+": "+result.Output)
			}
		}
	}
	return strings.Join(outputs, "; ")
}

// HealthHandler provides an HTTP handler that can be used to serve the health check endpoint.
func HealthHandler(healthchecker *HealthCheck) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Set("Content-Type", "application/health+json")

		result := healthchecker.Execute(ctx)

		// Map to HTTP response structure
		resp := HealthHttpResponse{
			ServiceID:   healthchecker.ServiceID,
			Description: healthchecker.Description,
			Version:     healthchecker.Version,
			ReleaseID:   healthchecker.ReleaseID,
			Status:      result.Status,
			Checks:      result.Checks,
			Output:      buildOutput(result.Checks),
		}

		if result.Status == checks.StatusFail {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		_ = json.NewEncoder(w).Encode(resp)
	}
}
