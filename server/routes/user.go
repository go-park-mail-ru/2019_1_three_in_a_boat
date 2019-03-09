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

func (h *UserHandler) Methods() map[string]struct{} {
	return map[string]struct{}{"PUT": {}, "GET": {}}
}

func (h *UserHandler) AuthRequired(method string) bool {
	if method == "PUT" {
		return true
	} else {
		return false
	}
}

func (h *UserHandler) CorsAllowed(method string) bool {
	return true
}
