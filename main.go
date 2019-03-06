package main

import (
	"flag"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/routes"
	"github.com/google/logger"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var verbose = flag.Bool("v", true, "print info level logs to stdout")
var logPath = flag.String("l", defaultLogPath, "path to the log file")
var sysLog = flag.Bool("sl", false, "log to syslog")

func main() {
	logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	logFile, err := os.OpenFile(
		*logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)

	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	defer logFile.Close()
	l := logger.Init("Hexagon Server", *verbose, *sysLog, logFile)
	defer l.Close()

	_db, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	mux := routes.GETRoutesMux(_db)
	s := makeServer(mux)

	logger.Fatal(s.ListenAndServe())
}
