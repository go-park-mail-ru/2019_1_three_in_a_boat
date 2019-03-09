// Package defines Server() which represents a hexagon API server
package server

// The file provides constants, structs and interfaces necessary for using the
// routes package

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/middleware"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/routes"
	"github.com/google/logger"
	"net/http"
	"sync"
	"time"
)

// Maps true/false into a string returned to the client in the status field
var RoutesMap = map[string]routes.Route{
	"/authors": {
		Handler: &routes.AuthorsHandler{},
		Name:    "authors",
	},
	"/users": {
		Handler: &routes.UsersHandler{},
		Name:    "users",
	},
	"/users/": {
		Handler: &routes.UserHandler{},
		Name:    "user",
	},
	"/signin": {
		Handler: &routes.SigninHandler{},
		Name:    "signin",
	},
}

var globalRouter = http.ServeMux{}
var routesMapOnce = sync.Once{}

// Singleton-like function, since router can be reused.
// Aggregates RoutesMap into a http.Handler, handling all acceptable requests
func GetRouter() http.Handler {
	routesMapOnce.Do(func() {
		var err error
		if err != nil {
			logger.Fatal("Failed to connect to DB")
		}

		logger.Info("Setting up router")
		for routeStr, routeObj := range RoutesMap {
			globalRouter.Handle(routeStr,
				middleware.Methods(
					middleware.CORS(
						middleware.Auth(
							routeObj.Handler, routeObj), routeObj), routeObj))
		}
		globalRouter.Handle("/", http.HandlerFunc(handlers.Handle404))
	})

	return &globalRouter
}

// Creates a new server with default settings and GetRouter() handler
func Server() *http.Server {
	return &http.Server{
		Addr:              ":3000",
		Handler:           GetRouter(),
		TLSConfig:         nil,
		ReadTimeout:       time.Second * 30,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 30,
		IdleTimeout:       time.Second * 60,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil, // purely useless, logging is handled elsewhere
	}
}
