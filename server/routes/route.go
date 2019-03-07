package routes

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Handler interface {
	http.Handler
}

type Route struct {
	Handler      Handler
	Methods      map[string]struct{}
	Middlewares  []mux.MiddlewareFunc
	AuthRequired bool
	CorsAllowed  bool
	Name         string
}
