package settings

const (
	dbUsername        = "hexagon"
	dbName            = "hexagon"
	dbHost            = "localhost"
	dbDefaultPassword = "ChangeBeforeDeploying"
)

// If -l is not specified, logs will be stored here
const DefaultLogPath = "logs/server.log"

// Path to file containing secret key. Meaningless when StoreKey is false
const SecretPath = "secret.rsa" // relative to the binary

// If true, the key will be loaded form SecretPath. If it does not exist, it
// will be created and a newly generated key will be written to it.
// If false, generates a new key on startup every time.
const StoreKey = true

// Simply the version returned to the client
const Version = "0.2"

// JSON values returned to the client, indicating whether the response was
// completed successfully. Is redundant, considering http status code, so
// provided just for convenience
var StatusMap = map[bool]string{
	true:  "ok",
	false: "error",
}

// Set-like map of allowed origins. If Origin belongs to this set, it will be
// returned to the client. Otherwise Access-Control remains unset.
var AllowedOrigins = map[string]struct{}{
	"http://localhost":               {},
	"https://three-in-a-boat.now.sh": {},
}
