// Package defines constants and singletons used across the entire server
// package. All the constants that affect the behavior of the server itself are
// listed in settings.go, constants that define client-related interaction
// specifics are listed in api.go
package server_settings

import (
	"flag"
	"os"
	"strconv"

	"github.com/google/logger"
	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats/pb"
	shared "github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/shared"
)

// If -l is not specified, logs will be stored here
const DefaultLogPath = "etc/logs/server.log"

var Verbose = flag.Bool("v", true, "print info level logs to stdout")
var LogPath = flag.String("l", DefaultLogPath, "path to the log file")
var SysLog = flag.Bool("sl", false, "log to syslog")
var ServerPort = flag.Int("p", shared.DefaultServerPort, "external API port")
var AuthAddress = flag.String(
	"auth",
	shared.DefaultAuthAddress+":"+strconv.Itoa(shared.DefaultAuthPort),
	"auth service address (with port)")

func SetUp() (*os.File, *logger.Logger, *grpc.ClientConn) {
	flag.Parse()
	logFile, l := shared.SetUpLog(*LogPath, *Verbose, *SysLog)

	// triggering the do.Once for logging and triggering fatal errors
	shared.DB()
	shared.GetAllowedOrigins()

	conn := shared.AuthConn(*AuthAddress)
	shared.AuthClient = pb.NewAuthClient(conn)
	return logFile, l, conn
}
