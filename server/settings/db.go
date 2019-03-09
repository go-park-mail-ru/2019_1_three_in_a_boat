package settings

import (
	"database/sql"
	"fmt"
	"github.com/google/logger"
	"os"
	"sync"
)

var setDbParamsOnce sync.Once

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

// makes a postgres connection string based on db* constants and PGPASSWORD
// environment variable
func makeConnStr() string {
	setDbParams()
	return fmt.Sprintf("postgres://%s:%s@%s/%s",
		dbUsername, dbPassword, dbHost, dbName)
}

var dbOnce = sync.Once{}
var db *sql.DB

// gets DB in a singleton-like manner
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
