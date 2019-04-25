// Package defines constants and singletons used across the entire server
// package. All the constants that affect the behavior of the server itself are
// listed in settings.go, constants that define client-related interaction
// specifics are listed in api.go
package server_settings

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/google/logger"
	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/formats/pb"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/settings/shared"
)

// If -l is not specified, logs will be stored here
const DefaultLogPath = "etc/logs/server.log"

var authConn *grpc.ClientConn

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
var ServerPort = flag.Int("p", DefaultServerPort, "external API port")
var AuthAddress = flag.String(
	"auth",
	DefaultAuthAddress+":"+strconv.Itoa(DefaultAuthPort),
	"external API port")

func SetUp() (*os.File, *logger.Logger, *grpc.ClientConn) {
	// using flag.Parse in init is discouraged so using this function which must
	// be called explicitly instead. Also logfiles need to be closed and we can't
	// return from init.

	flag.Parse()
	logFile, err := os.OpenFile(
		*LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	l := logger.Init("Hexagon Server", *Verbose, *SysLog, logFile)

	// triggering the do.Once for logging and triggering fatal errors
	DB()
	GetAllowedOrigins()

	return logFile, l, AuthConn()
}

var dialAuthOnce sync.Once

func AuthConn() *grpc.ClientConn {
	dialAuthOnce.Do(
		func() {
			var err error
			authConn, err = grpc.Dial(*AuthAddress, grpc.WithInsecure())
			if err != nil {
				logger.Fatalf("Failed to dial the auth server: %v", err)
			}
		})
	return authConn
}

var AuthClient = pb.NewAuthClient(AuthConn())
