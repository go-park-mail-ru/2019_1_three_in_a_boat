package game_settings

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
const DefaultLogPath = "etc/logs/game.log"

var Verbose = flag.Bool("v", true, "print info level logs to stdout")
var LogPath = flag.String("l", DefaultLogPath, "path to the log file")
var SysLog = flag.Bool("sl", false, "log to syslog")
var AuthAddress = flag.String(
	"auth",
	shared.DefaultAuthAddress+":"+strconv.Itoa(shared.DefaultAuthPort),
	"auth service address (with port)")
var GamePort = flag.Int("p", shared.DefaultGamePort, "game port")

func SetUp() (*os.File, *logger.Logger, *grpc.ClientConn) {
	// using flag.Parse in init is discouraged so using this function which must
	// be called explicitly instead. Also logfiles need to be closed and we can't
	// return from init.

	logFile, l := shared.SetUpLog(*LogPath, *Verbose, *SysLog)

	// triggering the do.Once for logging and triggering fatal errors
	shared.DB()
	shared.GetAllowedOrigins()
	conn := shared.AuthConn(*AuthAddress)
	shared.AuthClient = pb.NewAuthClient(conn)
	return logFile, l, conn
}
