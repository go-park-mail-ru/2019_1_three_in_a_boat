// Package defines Handlers for all roots used by the hexagon app
package routes

// The file provides constants, structs and interfaces necessary for using the
// routes package

import (
	"database/sql"
	"github.com/google/logger"
	"net/http"
	"sync"
)

// Maps true/false into a string returned to the client in the status field
var StatusMap = map[bool]string{
	true:  "ok",
	false: "error",
}

const Version = "0.1"

// Defines an HTTP handler on top of the http.Handler interface: adds the SetDB method
type Handler interface {
	http.Handler
	SetDB(*sql.DB)
	SetRoute(string)
	GetRoute() string
}

// The routes map. Values will be set up with db and route in a call to
// GETRoutesMap. Must not be used directly.
var getRoutesMap = map[string]Handler{
	"/authors": &AuthorsHandler{},
}

var getRoutesMapOnce = sync.Once{}

// Returns routes map in the form of {route: handler}. Only the first call is
// required to supply an actual sql.DB - subsequent calls can just pass nil.The
// function is goroutine-safe, however, the first succeeded call must provide a
// real sql.DB. Providing a nil in the first call will result in a fatal panic.
func GETRoutesMap(db *sql.DB) map[string]Handler {
	getRoutesMapOnce.Do(func() {
		if db == nil {
			logger.Fatal("GETRoutesMap called with nil the first time it was called")
		}

		for route, handler := range getRoutesMap {
			handler.SetDB(db)
			handler.SetRoute(route)
		}
	})

	return getRoutesMap
}

// same as the routesMap, but this one is a ServeMux based on the map.
var getRoutesMux = http.NewServeMux()
var getRoutesMuxOnce = sync.Once{}

// Returns an http.ServeMux based on GETRoutesMap's map. It calls GETRoutesMap,
// so the same concurrency rules apply.
func GETRoutesMux(db *sql.DB) *http.ServeMux {
	getRoutesMuxOnce.Do(func() {
		for route, handler := range GETRoutesMap(db) {
			getRoutesMux.Handle(route, handler)
		}
	})
	return getRoutesMux
}
