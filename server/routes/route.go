package routes

import (
	"net/http"
)

// Handler, as in http.Handler, but could be extended in the future
type Handler interface {
	http.Handler
	Methods() (methods map[string]struct{})
	AuthRequired(method string) bool
	CorsAllowed(method string) bool
}

// Represents all the data there's to know about a root. Every middleware checks
// if it should be applied to this particular route using this struct.
type Route struct {
	Handler Handler
	Name    string
}
