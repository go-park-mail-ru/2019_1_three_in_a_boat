package settings

import (
	"github.com/google/logger"
	"gopkg.in/square/go-jose.v2"
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
const UploadsPath = "media/images"

var ImageSize = [...]int{400, 400}

// Path to file containing secret key. Meaningless when StoreKey is false
const SecretPath = "secret.rsa" // relative to the binary

// If true, the key will be loaded form SecretPath. If it does not exist, it
// will be created and a newly generated key will be written to it.
// If false, generates a new key on startup every time.
const StoreKey = true

// Signing parameters, just don't change them and you'll be fine. Probably
const SigningAlgorithm = jose.RS256

var SignerOpts = jose.SignerOptions{
	ExtraHeaders: map[jose.HeaderKey]interface{}{
		"typ": "JWT",
	},
}

// Simply the version returned to the client
const Version = "0.9"

// JSON values returned to the client, indicating whether the response was
// completed successfully. Is redundant, considering http status code, so
// provided just for convenience
var StatusMap = map[bool]string{
	true:  "ok",
	false: "error",
}

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
