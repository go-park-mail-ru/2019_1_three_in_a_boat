// Package defines Server() which represents a hexagon API server
package main

// The file provides constants, structs and interfaces necessary for using the
// routes package

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/logger"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/middleware"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/routes"
)

// Maps URL paths into corresponding routes.Routes
var RoutesMap = map[string]routes.Handler{
	"/authors": &routes.AuthorsHandler{},
	"/users":   &routes.UsersHandler{},
	"/users/":  &routes.UserHandler{},
	"/signin":  &routes.SigninHandler{},
	"/":        &routes.CheckAuthHandler{},
	"/signout": &routes.SignOutHandler{},
	"/play":    &routes.SinglePlayHandler{},
}

var globalRouter = http.ServeMux{}
var routesMapOnce = sync.Once{}

// Singleton-like function, since router can be reused.
// Aggregates RoutesMap into a http.Handler, handling all acceptable requests
func GetRouter() http.Handler {
	routesMapOnce.Do(func() {
		logger.Info("Setting up router")
		for route, handler := range RoutesMap {
			globalRouter.Handle(route,
				middleware.Panic(
					middleware.Methods(
						middleware.CORS(
							middleware.CSRF(
								middleware.Auth(handler))))))
		}
	})

	return &globalRouter
}

// Creates a new server with default settings and GetRouter() handler
func Server(port int) *http.Server {
	return &http.Server{
		Addr:              ":" + strconv.Itoa(port),
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
