// Package defines constants and singletons used across the entire server
// package. All the constants that affect the behavior of the server itself are
// listed in settings.go, constants that define client-related interaction
// specifics are listed in api.go
package settings

import (
	_ "github.com/lib/pq"
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

const DefaultAuthAddress = "localhost"
const DefaultServerPort = 3000
const DefaultAuthPort = 3001
