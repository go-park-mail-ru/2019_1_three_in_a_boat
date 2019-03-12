package routes

import (
	"net/http"
)

// A handler that handles a ~multiple~ users resource. The handler itself is
// simply a struct that forwards the requests, depending on the method to one of
// GetUsers or PostUsers. Accepts PUT and POST requests. Implements
// route.Handler. For details on what do both resources do, see
// users_(method).go.
type UsersHandler struct{}

func (h *UsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		PostUsers(w, r)
	} else if r.Method == "GET" {
		GetUsers(w, r)
	}
}

func (h *UsersHandler) Settings() map[string]RouteSettings {
	return map[string]RouteSettings{
		"POST": {
			AuthRequired:           false,
			CorsAllowed:            true,
			CsrfProtectionRequired: true,
		},
		"GET": {
			AuthRequired:           false,
			CorsAllowed:            true,
			CsrfProtectionRequired: false,
		},
	}
}
