// Package defines Handlers for all roots used by the hexagon app
package routes

// The file provides constants, structs and interfaces necessary for using the
// routes package

import (
	"database/sql"
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

// Maps true/false into a string returned to the client in the status field
var StatusMap = map[bool]string{
	true:  "ok",
	false: "error",
}

const Version = "0.1"

var allowedOrigins = map[string]struct{}{
	"http://localhost": {}, "https://three-in-a-boat.now.sh": {}}

var routesMap = map[string]route{
	"/authors": {
		handler:      &AuthorsHandler{},
		methods:      map[string]struct{}{"GET": {}},
		middlewares:  []mux.MiddlewareFunc{},
		authRequired: false,
		corsAllowed:  true,
		name:         "authors",
	},
}

var globalRouter = mux.NewRouter()
var _db *sql.DB = nil

// Defines an HTTP handler on top of the http.Handler interface: adds the SetDB method
type Handler interface {
	http.Handler
}

type route struct {
	handler      Handler
	methods      map[string]struct{}
	middlewares  []mux.MiddlewareFunc
	authRequired bool
	corsAllowed  bool
	name         string
}

var routesMapOnce = sync.Once{}

func GetRouter(database *sql.DB) http.Handler {
	routesMapOnce.Do(func() {
		if database == nil {
			logger.Fatal("_db is nil in the first GetRouter call")
		}
		_db = database
		logger.Info("Setting up router")

		for routeStr, routeObj := range routesMap {
			globalRouter.Handle(routeStr,
				MethodMiddleware(
					CORSMiddleware(
						AuthMiddleware(
							routeObj.handler, routeObj), routeObj), routeObj)).
				Name(routeObj.name)
		}

		globalRouter.NotFoundHandler = http.HandlerFunc(Handle404)
	})

	return globalRouter
}
