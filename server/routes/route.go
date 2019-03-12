// Package provides handlers for all of the resources available in the app, as
// well as the struct, functions and interfaces needed to interact with them
// It extensively uses handlers package to provide consistent responses across
// all urls.
package routes

import (
	"net/http"
)

// Handler interface that embeds http.Handler. It provides additional methods,
// used by middlewares in order to understand how to handle a particular
// resource.
type Handler interface {
	http.Handler
	// Returns map of all allowed methods, mapping them to the RouteSettings
	Settings() (methods map[string]RouteSettings)
}

// Represents all the data there's to know about a root. Every middleware checks
// if it should be applied to this particular route using this struct, its
// Handler in particular.
type Route struct {
	Handler Handler
	Name    string
}
type RouteSettings struct {
	// Indicates whether unauthorized users are not allowed to work with the
	// resource. The authorization process is handled in the Auth middleware
	AuthRequired bool
	// Indicates whether cross-origin requests should/ be allowed. False indicates
	// that the resource is not supposed to be used by the client JS API.
	CorsAllowed bool
	// Indicates whether CSRF protection is required by the resource. The
	// protection is handled in the CSRF middleware
	CsrfProtectionRequired bool
}
