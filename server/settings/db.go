package settings

import (
	"database/sql"
	"fmt"
	"github.com/google/logger"
	"os"
	"sync"
)

var setDbParamsOnce sync.Once

// singleton-like function that sets the DB parameters from the environment
// or uses the default ones. Uses PGPASSWORD, PGUSERNAME, PGHOST, PGDBNAME
// environment variables.
func setDbParams() {
	setDbParamsOnce.Do(func() {
		if os.Getenv("PGPASSWORD") != "" {
			dbPassword = os.Getenv("PGPASSWORD")
		} else {
			logger.Warningln("Using debug password")
		}

		if os.Getenv("PGUSERNAME") != "" {
			dbUsername = os.Getenv("PGUSERNAME")
		}

		if os.Getenv("PGHOST") != "" {
			dbHost = os.Getenv("PGHOST")
		}

		if os.Getenv("PGDBNAME") != "" {
			dbName = os.Getenv("PGDBNAME")
		}
	})
}

// makes a postgres connection string based on setDbParams
func makeConnStr() string {
	setDbParams()
	return fmt.Sprintf("postgres://%s:%s@%s/%s",
		dbUsername, dbPassword, dbHost, dbName)
}

var dbOnce = sync.Once{}
var db *sql.DB

// Gets/creates DB in a singleton-like manner
func DB() *sql.DB {
	dbOnce.Do(func() {
		var err error
		logger.Info("Connecting to PostgreSQL...")
		db, err = sql.Open("postgres", makeConnStr())
		if err != nil {
			logger.Fatal("Failed to connect to DB: ", err)
		} else {
			err = db.Ping()
			if err != nil {
				logger.Fatal("Failed to ping DB: ", err)
			}
			logger.Info("Connected!")
		}
	})
	return db
}
