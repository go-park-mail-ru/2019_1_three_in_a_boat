package routes

import (
	"net/http"
)

type UsersHandler struct{}

func (h *UsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		PostUsers(w, r)
	} else if r.Method == "GET" {
		GetUsers(w, r)
	}
}

func (h *UsersHandler) Methods() map[string]struct{} {
	return map[string]struct{}{"POST": {}, "GET": {}}
}

func (h *UsersHandler) AuthRequired(method string) bool {
	if method == "PUT" {
		return true
	} else {
		return false
	}
}

func (h *UsersHandler) CorsAllowed(method string) bool {
	return true
}
