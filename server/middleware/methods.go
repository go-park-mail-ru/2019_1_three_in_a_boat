package middleware

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/routes"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/settings"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/utils"
	"net/http"
)

func MethodsMiddleware(next http.Handler, _route routes.Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, found := _route.Methods[r.Method]; found {
			next.ServeHTTP(w, r)
		} else if _route.CorsAllowed && r.Method == "OPTIONS" {
			CORSMiddleware(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					_, allowed := settings.AllowedOrigins[r.Header.Get("Origin")]
					if !allowed {
						utils.Handle405(w, r)
					} else {
						w.WriteHeader(http.StatusOK)
					}
				}), _route).ServeHTTP(w, r)
		} else {
			utils.Handle405(w, r)
		}
	})
}
