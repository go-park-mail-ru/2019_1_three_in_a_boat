// Package defines constants and singletons used across the entire server
// package. All the constants that affect the behavior of the server itself are
// listed in settings.go, constants that define client-related interaction
// specifics are listed in api.go
package auth_settings

import (
	"flag"
	"log"
	"os"

	"github.com/google/logger"
	"gopkg.in/square/go-jose.v2"

	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/settings/shared"
)

// If -l is not specified, logs will be stored here
const DefaultLogPath = "etc/logs/auth.log"

// Path to file containing secret key. Meaningless when StoreKey is false
const SecretPath = "secret.rsa" // relative to the binary

// If true, the key will be loaded form SecretPath. If it does not exist, it
// will be created and a newly generated key will be written to it.
// If false, generates a new key on startup every time.
const StoreKey = true

// Signing algorithm. The documentation recommends RS256 as the default one.
const SigningAlgorithm = jose.RS256

var Verbose = flag.Bool("v", true, "print info level logs to stdout")
var LogPath = flag.String("l", DefaultLogPath, "path to the log file")
var SysLog = flag.Bool("sl", false, "log to syslog")
var AuthPort = flag.Int("p", DefaultAuthPort, "auth service port (must be shared across all services)")

func SetUp() (*os.File, *logger.Logger) {
	// using flag.Parse in init is discouraged so using this function which must
	// be called explicitly instead

	flag.Parse()
	logFile, err := os.OpenFile(
		*LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	l := logger.Init("Hexagon Server", *Verbose, *SysLog, logFile)

	// triggering the do.Once for logging and triggering fatal errors
	GetSigner()
	DB()

	return logFile, l
}
