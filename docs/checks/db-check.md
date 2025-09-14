# DB Check

The DB Check provides checks to monitor the health of a SQL database connection. It can perform a simple ping to verify connectivity and monitor database connection pool metrics.

## Ping Check

Ping Check verifies that the database is reachable by performing a ping operation.

### Configuration Options

The Ping Check can be configured using the following options:

- `WithPingName(name string)`: Sets the name of the check.
- `WithPingDB(db DatabasePinger)`: Sets the database connection to be used for the check.
- `WithPingTimeout(timeout time.Duration)`: Sets the timeout for the ping operation (default is 5 seconds).

### Example Usage

```go
package main

import (
  "database/sql"
  "net/http"
  "time"

  _ "github.com/lib/pq" // Import the PostgreSQL driver

  "github.com/brpaz/go-healthcheck/v2/checks/dbcheck"
)

func main() {
  // Initialize the database connection
  db, err := sql.Open("postgres", "user=youruser dbname=yourdb sslmode=disable")
  if err != nil {
    panic(err)
  }
  defer db.Close()

  // Create a new Ping Check
  dbPingCheck := dbcheck.NewPing(
    dbcheck.WithPingName("postgres-ping"),
    dbcheck.WithPingDB(db),
    dbcheck.WithPingTimeout(2*time.Second),
  )
}
```

## Connections Check

Connections Check monitors the number of open connections in the database connection pool and compares it against defined thresholds.

### Configuration Options

The Connections Check can be configured using the following options:

- `WithConnectionsName(name string)`: Sets the name of the check.
- `WithConnectionsDB(db DatabaseStatsProvider)`: Sets the database connection to be used for the check.
- `WithConnectionsTimeout(timeout time.Duration)`: Sets the timeout for the check operation (default is 5 seconds).
- `WithConnectionsWarnThreshold(threshold float64)`: Sets the warning threshold as a percentage (0-100) of max connections (default is 80.0).
- `WithConnectionsFailThreshold(threshold float64)`: Sets the failure threshold as a percentage (0-100) of max connections (default is 100.0).

### Example Usage

```go
package main

import (
  "database/sql"
  "net/http"
  "time"

  _ "github.com/lib/pq" // Import the PostgreSQL driver

  "github.com/brpaz/go-healthcheck/v2/checks/dbcheck"
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
  dbConnectionsCheck := dbcheck.NewConnections(
    dbcheck.WithConnectionsName("postgres-connections"),
    dbcheck.WithConnectionsDB(db),
    dbcheck.WithConnectionsWarnThreshold(80.0),
    dbcheck.WithConnectionsFailThreshold(95.0),
  )
}
```

