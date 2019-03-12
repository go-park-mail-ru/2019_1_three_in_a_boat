package routes

import "net/http"

// A handler that handles a ~single~ user resource. The handler itself is simply
// a struct that forwards the requests, depending on the method to one of
// PutUser or GetUser. Accepts PUT and GET requests. Implements route.Handler.
// For details on what do both resources do, see user_(method).go.
type UserHandler struct{}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		PutUser(w, r)
	} else if r.Method == "GET" {
		GetUser(w, r)
	}
}

func (h *UserHandler) Settings() map[string]RouteSettings {
	return map[string]RouteSettings{
		"PUT": {
			AuthRequired:           true,
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