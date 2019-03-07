// Package defines GetRouter function which returns mux ready to be plugged
// into http.ServeHttp
package server

// The file provides constants, structs and interfaces necessary for using the
// routes package

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/middleware"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/routes"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/utils"
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
	"time"
)

// Maps true/false into a string returned to the client in the status field

var routesMap = map[string]routes.Route{
	"/authors": {
		Handler:      &routes.AuthorsHandler{},
		Methods:      map[string]struct{}{"GET": {}},
		Middlewares:  []mux.MiddlewareFunc{},
		AuthRequired: false,
		CorsAllowed:  true,
		Name:         "authors",
	},
}

var globalRouter = mux.NewRouter()

// Defines an HTTP handler on top of the http.Handler interface: adds the SetDB method

var routesMapOnce = sync.Once{}

func GetRouter() http.Handler {
	routesMapOnce.Do(func() {
		var err error
		if err != nil {
			logger.Fatal("Failed to connect to DB")
		}

		logger.Info("Setting up router")
		for routeStr, routeObj := range routesMap {
			globalRouter.Handle(routeStr,
				middleware.MethodsMiddleware(
					middleware.CORSMiddleware(
						middleware.AuthMiddleware(
							routeObj.Handler, routeObj), routeObj), routeObj)).
				Name(routeObj.Name)
		}
		globalRouter.NotFoundHandler = http.HandlerFunc(utils.Handle404)
	})

	return globalRouter
}

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
