# Getting Started

## Installation

To install the library, use go get:

```bash
go get github.com/brpaz/go-healthcheck
```

## Basic Usage

The library provides a healthcheck service and an HTTP handler to expose the healthcheck endpoint.

```go
package main

import (
    "github.com/brpaz/go-healthcheck"
    "github.com/brpaz/go-healthcheck/checks/httpcheck"
    "net/http"
)

func main() {
    hc := healthcheck.New(
        healthcheck.WithServiceName("my-service"),
        healthcheck.WithDescription("My Service Healthcheck"),
        healthcheck.WithVersion("1.0.0"),
        healthcheck.WithReleaseID("sha256:abcdef1234567890"),
        healthcheck.WithChecks(
            httpcheck.New(
                httpcheck.WithName("http:google"),
                httpcheck.WithURL("https://www.google.com"),
                httpcheck.WithExpectedStatus(200)
            ),
        ),
    )
    http.Handle("/health", healthcheck.HealthHandler(hc))
    http.ListenAndServe(":8080", nil)
}
```

When requesting the `/health` endpoint, you will receive a JSON response similar to:

```json
{
  "status": "pass",
  "service": "my-service",
  "description": "My Service Healthcheck",
  "version": "1.0.0",
  "releaseId": "sha256:abcdef1234567890",
  "checks": {
    "http:google": {
      "status": "pass",
      "observedValue": 5,
      "observedUnit": "ms",
      "time": "2024-06-01T12:00:00Z"
    }
  ]
}
```

## Available Checks

To check the specific checks documention, please refer to the [checks documentation](./checks/index.md).
