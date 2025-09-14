## SQL Check

The SQL Check verifies that a SQL database is reachable. This is useful for monitoring the availability of SQL databases such as MySQL, PostgreSQL, SQLite, and others.

## Configuration

The SQL Check can be configured using the following options:

- `WithName(name string)`: Sets the name of the check.
- `WithDB(db *sql.DB)`: Sets the database connection to be used for the check.
- `WithTimeout(timeout time.Duration)`: Sets the timeout for the ping operation (default is 5 seconds).
- `WithMetrics(includeMetrics bool)`: Enables collection of database connection pool metrics (default is false).

## Connection Pool Metrics

When metrics are enabled with `WithMetrics(true)`, the health check will return separate sub-checks for each connection pool metric, allowing for granular monitoring and alerting:

- **`sql-check:open-connections`**: Current number of established connections
- **`sql-check:in-use-connections`**: Number of connections currently being used
- **`sql-check:idle-connections`**: Number of idle connections available
- **`sql-check:max-open-connections`**: Maximum number of open connections allowed
- **`sql-check:wait-count`**: Total number of times a connection was waited for
- **`sql-check:wait-duration`**: Total time spent waiting for connections

Each metric has its own `ObservedValue` and `ObservedUnit`, making it easy to set up monitoring dashboards and alerts on specific thresholds (e.g., alert when idle connections drop below 2, or when wait count increases rapidly).

**Note:** Connection count metrics use empty `observedUnit` since counts are dimensionless. Only time-based metrics like `wait-duration` and `sql-check` (ping time) use units like "ms".

## Example

```go
package main

import (
    "database/sql"
    "time"
    "github.com/brpaz/go-healthcheck"
    "github.com/brpaz/go-healthcheck/checks/sqlcheck"
    _ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

func main() {
    // Open a database connection
    db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Create a basic health check
    basicCheck := sqlcheck.New(
        sqlcheck.WithDB(db),
        sqlcheck.WithTimeout(2 * time.Second),
    )

    // Create a health check with metrics enabled
    metricsCheck := sqlcheck.New(
        sqlcheck.WithDB(db),
        sqlcheck.WithTimeout(2 * time.Second),
        sqlcheck.WithMetrics(true), // Enable connection pool metrics
    )

    // Use the checks...
}
```

## Sample Output

**Basic check output (1 result):**
```json
{
  "componentId": "sql-check",
  "status": "pass",
  "output": "",
  "observedValue": 1,
  "observedUnit": "ms"
}
```

**With metrics enabled (7 results):**
```json
[
  {
    "componentId": "sql-check",
    "status": "pass",
    "output": "",
    "observedValue": 1,
    "observedUnit": "ms"
  },
  {
    "componentId": "sql-check:open-connections",
    "status": "pass",
    "output": "",
    "observedValue": 5,
    "observedUnit": ""
  },
  {
    "componentId": "sql-check:in-use-connections",
    "status": "pass",
    "output": "",
    "observedValue": 2,
    "observedUnit": ""
  },
  {
    "componentId": "sql-check:idle-connections",
    "status": "pass",
    "output": "",
    "observedValue": 3,
    "observedUnit": ""
  },
  {
    "componentId": "sql-check:max-open-connections",
    "status": "pass",
    "output": "",
    "observedValue": 25,
    "observedUnit": ""
  },
  {
    "componentId": "sql-check:wait-count",
    "status": "pass",
    "output": "",
    "observedValue": 10,
    "observedUnit": ""
  },
  {
    "componentId": "sql-check:wait-duration",
    "status": "pass",
    "output": "",
    "observedValue": 15,
    "observedUnit": "ms"
  }
]
```

## Monitoring Use Cases

The separate metric sub-checks enable advanced monitoring scenarios:

**Alerting Examples:**
- Alert when `idle-connections` < 2 (connection pool exhaustion risk)
- Alert when `wait-count` increases by > 100 over 5 minutes (connection contention)
- Alert when `wait-duration` > 1000ms (slow connection acquisition)
- Alert when `in-use-connections` / `max-open-connections` > 0.8 (high utilization)

**Dashboard Metrics:**
- Graph `open-connections` over time to see connection usage patterns
- Monitor `wait-duration` trends to identify performance degradation
- Track `idle-connections` to optimize pool sizing
