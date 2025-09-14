## DB Check

The DB Check provides checks to monitor the health of a SQL database connection. It can perform a simple ping to verify connectivity and monitor database connection pool metrics.

## Available Checks

### Ping Check

The Ping Check verifies that a connection to the database can be established by performing a ping operation. This is useful for monitoring the general availability of the database.

### Connections Check

The Connections Check monitors the number of open database connections against configured thresholds. It automatically detects the maximum connection limit from the database settings.

## Configuration

### Ping Check

The Ping Check can be configured using the following options:

- `WithPingName(name string)`: Sets the name of the check.
- `WithPingDB(db DatabasePinger)`: Sets the database connection to be used for the check.
- `WithPingTimeout(timeout time.Duration)`: Sets the timeout for the ping operation (default is 5 seconds).

### Connections Check

The Connections Check can be configured using the following options:

- `WithName(name string)`: Sets the name of the check.
- `WithDB(db DatabaseStatsProvider)`: Sets the database connection to be used for the check.
- `WithTimeout(timeout time.Duration)`: Sets the timeout for the check operation (default is 5 seconds).
- `WithWarnThreshold(threshold float64)`: Sets the warning threshold as a percentage (0.0-1.0) of max connections (default is 0.8).
- `WithFailThreshold(threshold float64)`: Sets the failure threshold as a percentage (0.0-1.0) of max connections (default is 0.9).

### Example Usage

```go
package main

import (
  "database/sql"
  "net/http"
  "time"

  _ "github.com/lib/pq" // Import the PostgreSQL driver

  "github.com/brpaz/go-healthcheck/checks/dbcheck/connectionscheck"
)

func main() {
  // Initialize the database connection
  db, err := sql.Open("postgres", "user=youruser dbname=yourdb sslmode=disable")
  if err != nil {
    panic(err)
  }
  defer db.Close()

  // Configure connection pool
  db.SetMaxOpenConns(100)
  db.SetMaxIdleConns(10)

  // Create a new Connections Check
  dbConnectionsCheck := connectionscheck.New(
    connectionscheck.WithName("postgres-connections"),
    connectionscheck.WithDB(db),
    connectionscheck.WithWarnThreshold(0.8),
    connectionscheck.WithFailThreshold(0.95),
  )

  // Create health checker with both checks
  checker := healthcheck.New(
    healthcheck.WithChecks(dbConnectionsCheck),
  )
}
```

