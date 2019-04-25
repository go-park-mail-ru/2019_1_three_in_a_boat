// Package defines constants and singletons used across the entire server
// package. All the constants that affect the behavior of the server itself are
// listed in settings.go, constants that define client-related interaction
// specifics are listed in api.go
package settings

import (
	"flag"
	"github.com/google/logger"
	"gopkg.in/square/go-jose.v2"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	dbUsername = "hexagon"
	dbName     = "hexagon"
	dbHost     = "localhost"
	dbPassword = "ChangeBeforeDeploying"
)

// If -l is not specified, logs will be stored here
const DefaultLogPath = "logs/server.log"

const AuthPort = 8000

var ImageSize = [...]int{400, 400}

// Path to file containing secret key. Meaningless when StoreKey is false
const SecretPath = "secret.rsa" // relative to the binary

// If true, the key will be loaded form SecretPath. If it does not exist, it
// will be created and a newly generated key will be written to it.
// If false, generates a new key on startup every time.
const StoreKey = true

// Signing algorithm. The documentation recommends RS256 as the default one.
const SigningAlgorithm = jose.RS256

// CSRF when debugging is annoying af, hence this setting
const EnableCSRF = true

// Simply the version returned to the client
const Version = "0.9"

// Set-like map of allowed origins. If Origin belongs to this set, it will be
// returned to the client. Otherwise Access-Control remains unset.
var allowedOrigins = map[string]struct{}{
	"http://localhost":               {},
	"http://localhost:8080":          {},
	"http://localhost:3000":          {},
	"http://127.0.0.1":               {},
	"http://127.0.0.1:8080":          {},
	"http://127.0.0.1:3000":          {},
	"https://three-in-a-boat.now.sh": {},
}

var allowedOriginsOnce sync.Once

// concats allowedOrigins with everything found in ALLOWED_ORIGINS environment
// variable (extracts urls split by ;)
func GetAllowedOrigins() map[string]struct{} {
	allowedOriginsOnce.Do(func() {
		origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ";")
		for _, origin := range origins {
			if origin != "" {
				allowedOrigins[strings.TrimSpace(origin)] = struct{}{}
			}
		}

		logger.Info("Listing allowed origins:")
		for origin := range allowedOrigins {
			logger.Info("\t", origin)
		}
	})

	return allowedOrigins
}

var Verbose = flag.Bool("v", true, "print info level logs to stdout")
var LogPath = flag.String("l", DefaultLogPath, "path to the log file")
var SysLog = flag.Bool("sl", false, "log to syslog")
var ServerPort = flag.Int("p", 3000, "port to listen at")

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
	GetAllowedOrigins()

	return logFile, l
}
