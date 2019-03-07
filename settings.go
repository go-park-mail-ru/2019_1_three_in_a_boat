package main

import (
	"database/sql"
	"fmt"
	"github.com/google/logger"
	"net/http"
	"os"
	"time"
)

const (
	dbUsername        = "hexagon"
	dbName            = "hexagon"
	dbHost            = "localhost"
	dbDefaultPassword = "ChangeBeforeDeploying"
)

const (
	defaultLogPath = "logs/server.log"
)

func makeConnStr() string {
	dbPassword := os.Getenv("PGPASSWORD")
	if dbPassword == "" {
		dbPassword = dbDefaultPassword
		logger.Warningln("Using debug password")
	}

	return fmt.Sprintf("postgres://%s:%s@%s/%s",
		dbUsername, dbPassword, dbHost, dbName)
}

func connectToDB() (*sql.DB, error) {
	return sql.Open("postgres", makeConnStr())
}

func makeServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              ":3000",
		Handler:           handler,
		TLSConfig:         nil,
		ReadTimeout:       time.Second * 30,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 30,
		IdleTimeout:       time.Second * 60,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil, // purely useless, logging is handled elsewhere
	}
}
