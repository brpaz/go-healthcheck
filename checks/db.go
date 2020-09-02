package checks

import (
	"database/sql"
	"fmt"

	"github.com/brpaz/go-healthcheck"
)

type DBChecker struct {
	componentID string
	db          *sql.DB
}

// NewDBChecker initializes an instance of the Database Checker.
func NewDBChecker(componentID string, db *sql.DB) DBChecker {
	return DBChecker{
		componentID: componentID,
		db:          db}
}

// Execute Executes the check.
// This check tries to do a simple "SELECT 1" using the "db" object passed to the checker.
// The check will be considered "passing" if the query ran successfully.
func (c DBChecker) Execute() map[string][]healthcheck.Check {

	checkName := fmt.Sprintf("%s:status", c.componentID)

	status := healthcheck.Pass
	output := ""

	var ok string

	if err := c.db.QueryRow("SELECT 1").Scan(&ok); err != nil {
		status = healthcheck.Fail
		output = err.Error()
	}

	return map[string][]healthcheck.Check{
		checkName: {
			{
				ComponentType: componentTypeDatastore,
				ComponentID:   c.componentID,
				Status:        status,
				Output:        output,
			},
		},
	}
}
