package settings

import (
	"database/sql"
	"fmt"
	"github.com/google/logger"
	"os"
	"sync"
)

// makes a postgres connection string based on db* constants and PGPASSWORD
// environment variable
func makeConnStr() string {
	dbPassword := os.Getenv("PGPASSWORD")
	if dbPassword == "" {
		dbPassword = dbDefaultPassword
		logger.Warningln("Using debug password")
	}

	return fmt.Sprintf("postgres://%s:%s@%s/%s",
		dbUsername, dbPassword, dbHost, dbName)
}

var dbOnce = sync.Once{}
var db *sql.DB

// gets DB in a singleton-like manner
func DB() *sql.DB {
	dbOnce.Do(func() {
		var err error
		db, err = sql.Open("postgres", makeConnStr())
		if err != nil {
			logger.Fatal("Failed to connect to DB")
		}
	})
	return db
}
