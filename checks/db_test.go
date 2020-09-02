package checks_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brpaz/go-healthcheck"
	"github.com/brpaz/go-healthcheck/checks"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestDbChecker_Success(t *testing.T) {

	dbConn, mock, err := sqlmock.New()
	mock.ExpectQuery("SELECT 1").WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow("1"))

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbConn.Close()

	if err != nil {
		t.Fatal(err)
	}

	checker := checks.NewDBChecker("main-db", dbConn)

	result := checker.Execute()

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "main-db", result["main-db:status"][0].ComponentID)
	assert.Equal(t, healthcheck.Pass, result["main-db:status"][0].Status)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDbChecker_Fail(t *testing.T) {

	dbConn, err := sql.Open("mysql", "root@/blog")

	if err != nil {
		t.Fatal(err)
	}

	checker := checks.NewDBChecker("main-db", dbConn)

	result := checker.Execute()

	resultItem := result["main-db:status"][0]

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "main-db", resultItem.ComponentID)
	assert.Equal(t, healthcheck.Fail, resultItem.Status)
	assert.NotEmpty(t, resultItem.Output)
}
