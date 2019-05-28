// Package defines constants and singletons used across the entire server
// package. All the constants that affect the behavior of the server itself are
// listed in settings.go, constants that define client-related interaction
// specifics are listed in api.go
package settings

import (
	"flag"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/google/logger"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats/pb"
)

var (
	dbUsername = "hexagon"
	dbName     = "hexagon"
	dbHost     = "localhost"
	dbPassword = "ChangeBeforeDeploying"
)

// Simply the version returned to the client
const Version = "0.95"

// Specifies the string returned to the client when the entry in the DB is NULL or empty
const DefaultImgName = "default.png"

const DefaultServerPort = 3000

const DefaultAuthAddress = "localhost"
const DefaultAuthPort = 3001
const DefaultGamePort = 3002
const DefaultChatPort = 3003

// Lifespan of a JWE Auth token - in days
const JWETokenLifespan = 30

// Regulates the length of a CSRF Token (in bytes). 20 is probably ok.
const CSRFTokenLength = 20

// Lifespan of a CSRFToken - in days
const CSRFTokenLifespan = 7

// CSRF when debugging is annoying af, hence this setting. False = no checking.
const EnableCSRF = true

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

var authConn *grpc.ClientConn
var dialAuthOnce sync.Once

func AuthConn(addr string) *grpc.ClientConn {
	dialAuthOnce.Do(
		func() {
			var err error
			authConn, err = grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				logger.Fatalf("Failed to dial the auth server: %v", err)
			}
		})
	return authConn
}

var AuthClient pb.AuthClient

func SetUpLog(path string, verbose, sysLog bool) (*os.File, *logger.Logger) {
	// using flag.Parse in init is discouraged so using this function which must
	// be called explicitly instead. Also logfiles need to be closed and we can't
	// return from init.

	flag.Parse()
	logFile, err := os.OpenFile(
		path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return logFile, logger.Init("Hexagon Server", verbose, sysLog, logFile)
}
