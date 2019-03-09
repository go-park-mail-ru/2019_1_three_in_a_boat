package routes

import (
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"net/http"
)

type SignOutHandler struct{}

func (h *SignOutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Unauthorize(w, r)
	Handle200(w, r, nil)
}

func (h *SignOutHandler) Methods() map[string]struct{} {
	return map[string]struct{}{"GET": {}}
}

func (h *SignOutHandler) AuthRequired(method string) bool {
	return false
}

func (h *SignOutHandler) CorsAllowed(method string) bool {
	return true
}
