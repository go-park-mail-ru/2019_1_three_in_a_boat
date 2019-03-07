package settings

const (
	dbUsername        = "hexagon"
	dbName            = "hexagon"
	dbHost            = "localhost"
	dbDefaultPassword = "ChangeBeforeDeploying"
)

const DefaultLogPath = "logs/server.log"

const Version = "0.2"

var StatusMap = map[bool]string{
	true:  "ok",
	false: "error",
}

// Set-like map of allowed origins
var AllowedOrigins = map[string]struct{}{
	"http://localhost":               {},
	"https://three-in-a-boat.now.sh": {},
}
