// Package defines middleware working with routes.Route
package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/routes"
)

// CORS-preflight-friendly Methods middleware based on routes.Route and
// routes.Handler. Uses _route.Handler.Settings() to determine what methods are
// allowed. OPTIONS request to handle a potential preflight is allowed if it
// contains Access-Control-Request-Method and that method allows CORS. Simply
// returns the headers and does not call the handler in that case.
func Methods(next routes.Handler) routes.Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, found := next.Settings()[r.Method]; found {
			next.ServeHTTP(w, r)
		} else if r.Method == "OPTIONS" && next.Settings()[r.Header.Get(
			"Access-Control-Request-Method")].CorsAllowed {
			CORS(HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				},
				next.Settings())).ServeHTTP(w, r)
		} else {
			handlers.Handle405(w, r)
		}
	}, next.Settings())
}
