package chat_settings

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	"github.com/google/logger"
)

var (
	dbUsername = "hexagon"
	dbName     = "hexagon"
	dbHost     = "localhost"
	dbPassword = "ChangeBeforeDeploying"
)

var setDbParamsOnce sync.Once

// singleton-like function that sets the DB parameters from the environment
// or uses the default ones. Uses PGPASSWORD, PGUSERNAME, PGHOST, PGDBNAME
// environment variables. The parameters passed to the function are only
// used if they're non-empty.
func SetDbParams(pwd, username, host, name string) {
	setDbParamsOnce.Do(func() {
		if pwd != "" {
			dbPassword = pwd
		} else if os.Getenv("CHAT_PGPASSWORD") != "" {
			dbPassword = os.Getenv("CHAT_PGPASSWORD")
		} else {
			logger.Warningln("Using debug password")
		}

		if username != "" {
			dbUsername = username
		} else if os.Getenv("PGUSERNAME") != "" {
			dbUsername = os.Getenv("PGUSERNAME")
		}

		if host != "" {
			dbHost = host
		} else if os.Getenv("PGHOST") != "" {
			dbHost = os.Getenv("PGHOST")
		}

		if name != "" {
			dbName = name
		} else if os.Getenv("PGDBNAME") != "" {
			dbName = os.Getenv("PGDBNAME")
		}
	})
}

// makes a postgres connection string based on SetDbParams
func makeConnStr() string {
	SetDbParams("", "", "", "")
	return fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable",
		dbUsername, dbPassword, dbHost, dbName)
}

var dbOnce = sync.Once{}
var __db *sql.DB // don't ask

// Gets/creates sql.DB in a singleton-like manner using makeConnStr and
// SetDbParams with empty parameters
func DB() *sql.DB {
	dbOnce.Do(func() {
		var err error
		logger.Info("Connecting to PostgreSQL...")
		__db, err = sql.Open("postgres", makeConnStr())
		if err != nil {
			logger.Fatal("Failed to connect to DB: ", err)
		} else {
			err = __db.Ping()
			if err != nil {
				logger.Fatal("Failed to ping DB: ", err)
			}
			logger.Info("Connected!")
		}
	})
	return __db
}
