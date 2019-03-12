// Package defines middleware working with routes.Route
package middleware

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/routes"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/settings"
	"net/http"
)

// CORS-preflight-friendly Methods middleware based on routes.Route and
// routes.Handler. Uses _route.Handler.Settings() to determine what methods are
// allowed. OPTIONS request to handle a potential preflight is allowed if it
// contains Access-Control-Request-Method and that method allows CORS. Simply
// returns the headers and does not call the handler in that case.
func Methods(next http.Handler, _route routes.Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, found := _route.Handler.Settings()[r.Method]; found {
			next.ServeHTTP(w, r)
		} else if r.Method == "OPTIONS" && _route.Handler.Settings(
			    )[r.Header.Get("Access-Control-Request-Method")].CorsAllowed {
			CORS(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					_, allowed := settings.GetAllowedOrigins()[r.Header.Get("Origin")]
					if !allowed {
						handlers.Handle405(w, r)
					} else {
						w.WriteHeader(http.StatusOK)
					}
				}), _route).ServeHTTP(w, r)
		} else {
			handlers.Handle405(w, r)
		}
	})
}
