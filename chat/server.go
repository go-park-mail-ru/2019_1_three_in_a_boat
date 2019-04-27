// Package defines Server() which represents a hexagon API server
package main

// The file provides constants, structs and interfaces necessary for using the
// routes package. It contains a little bit of copypasta but there isn't really
// any logic in it so I think it's fine.

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/logger"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils/middleware"
)

// Maps URL paths into corresponding routes.Routes
var RoutesMap = map[string]http_utils.Handler{
	"/chat":          &ChatHandler{},
	"/chat/paginate": &ChatPaginationHandler{},
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
								middleware.Auth(handler)))))) // ))0)
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
