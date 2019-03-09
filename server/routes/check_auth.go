package routes

import (
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"net/http"
)

type CheckAuthHandler struct{}

// Dummy empty resource, to provide the user with the UserData.
// All the checking is handled in the Auth middleware, so we just need a cheap
// entrypoint.
func (h *CheckAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "" && r.URL.Path != "/" {
		Handle404(w, r)
	} else {
		Handle200(w, r, nil)
	}
}

func (h *CheckAuthHandler) Methods() map[string]struct{} {
	return map[string]struct{}{"GET": {}}
}

func (h *CheckAuthHandler) AuthRequired(method string) bool {
	return false
}

func (h *CheckAuthHandler) CorsAllowed(method string) bool {
	return true
}
