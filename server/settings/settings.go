package settings

import (
	"database/sql"
	"fmt"
	"github.com/google/logger"
	"os"
	"sync"
)

const (
	dbUsername        = "hexagon"
	dbName            = "hexagon"
	dbHost            = "localhost"
	dbDefaultPassword = "ChangeBeforeDeploying"
)

const DefaultLogPath = "logs/server.log"

const Version = "0.1"

var StatusMap = map[bool]string{
	true:  "ok",
	false: "error",
}

var AllowedOrigins = map[string]struct{}{
	"http://localhost":               {},
	"https://three-in-a-boat.now.sh": {},
}

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
