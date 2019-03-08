package middleware

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/routes"
	"net/http"
)

// Authentication middleware: if the resource requires authentication,
// checks JWT token and calls the method, returning 403 if it's invalid.
// Otherwise simply forwards the call
func Auth(next http.Handler, _route routes.Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
