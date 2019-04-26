package routes

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils/handlers"
)

// Handles Signout resource. Only accepts GET requests. Implements
// routes.Handler interface, which extends http.Handler. All it does is removing
// the auth cookie using handlers.Unauthorize and returns 200. The data in case
// of a successful response (and in case of an error, really) will be nil.
type SignOutHandler struct{}

func (h *SignOutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	Unauthorize(w)
	Handle200(w, r, nil)
}

func (h *SignOutHandler) Settings() map[string]http_utils.RouteSettings {
	return map[string]http_utils.RouteSettings{
		"POST": {
			AuthRequired:           false,
			CorsAllowed:            true,
			CsrfProtectionRequired: false,
		},
	}
}
