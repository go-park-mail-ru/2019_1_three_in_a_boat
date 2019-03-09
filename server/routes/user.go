package routes

import "net/http"

type UserHandler struct{}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// PutUser(w, r)
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
