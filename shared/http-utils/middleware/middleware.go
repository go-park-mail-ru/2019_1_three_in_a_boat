// The package defines middlewares used by the application. Middlewares have
// signature func(routes.Handler) routes.Handler, which are compatible with
// generic func(http.Handler) http.Handler middlewares, but also carry the
// settings of the route to the nested middlewares, because settings are used
// pretty much in every one of them.
package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils"
)

// Generic handler used by HandlerFunc to create a routes.Handler-compatible
// handler to pass between middlewares.
type GenericHandler struct {
	settings   map[string]http_utils.RouteSettings
	handleFunc func(w http.ResponseWriter, r *http.Request)
}

func (gh *GenericHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gh.handleFunc(w, r)
}

func (gh *GenericHandler) Settings() map[string]http_utils.RouteSettings {
	return gh.settings
}

func HandlerFunc(handler func(w http.ResponseWriter, r *http.Request),
	settings map[string]http_utils.RouteSettings) http_utils.Handler {
	return &GenericHandler{settings: settings, handleFunc: handler}
}
