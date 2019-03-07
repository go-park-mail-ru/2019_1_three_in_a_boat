package routes

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"net/http"
)

type UsersHandler struct{}

func (h *UsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		handlers.Handle405(w, r)
	} else {
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
