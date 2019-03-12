package middleware

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/routes"
	"net/http"
)

// CORS Middleware: adds Access-Control headers if request's Origin is allowed
// See settings for the allowed origins.
func CSRF(next http.Handler, _route routes.Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if the protection is required - check for it right away
		if _route.Handler.Settings()[r.Method].CsrfProtectionRequired {
			cookieToken, err := r.Cookie("csrf") // err ErrNoCookie only
			headerToken := r.Header.Get("X-CSRF-Token")
			if err != nil || headerToken != cookieToken.Value {
				handlers.LogError(
					0, "CSRF Validation Failed", r)
				handlers.Handle403Msg(w, r, formats.ErrCSRF)
				return
			}
		}
		// if it's not - set the token, if not already set
		cookie, err := r.Cookie("csrf")
		if err != nil || cookie.Value == "" {
			err = handlers.SetNewCsrfToken(w)
			if err != nil {
				handlers.LogError(
					0, fmt.Sprint("Failed to set CSRF Token: ", err), r)
				// still continue execution though
			}
		}
		next.ServeHTTP(w, r)
	})
}
