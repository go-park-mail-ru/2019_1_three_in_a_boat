package middleware

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/routes"
	"net/http"
)

func AuthMiddleware(next http.Handler, _route routes.Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
