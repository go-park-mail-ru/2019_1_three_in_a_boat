package routes

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils/handlers"
)

// Provides a convenient way for the frontend to check whether the user is
// authorized. The handler itself does nothing, but the auth middleware will
// return the user data in the response. It also handles 404 requests and is
// expected to be bound to the "/" path. Returns nothing in the JSON response
// data. There's probably a better way..
type CheckAuthHandler struct{}

// Handles the GET request.
func (h *CheckAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	if r.URL.Path != "" && r.URL.Path != "/" {
		Handle404(w, r)
	} else {
		Handle200(w, r, nil)
	}
}

func (h *CheckAuthHandler) Settings() map[string]http_utils.RouteSettings {
	return map[string]http_utils.RouteSettings{
		"GET": {
			AuthRequired:           false,
			CorsAllowed:            true,
			CsrfProtectionRequired: false,
		},
	}
}
