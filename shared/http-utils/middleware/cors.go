package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils"
	settings "github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/shared"
)

// CORS Middleware: adds Access-Control headers if request's Origin is allowed
// See settings for the allowed origins. Handles OPTIONS requests.
func CORS(next http_utils.Handler) http_utils.Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		_, allowed := settings.GetAllowedOrigins()[origin]
		method := r.Method
		if method == "OPTIONS" {
			method = r.Header.Get("Access-Control-Request-Method")
		}

		if allowed && next.Settings()[method].CorsAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers",
				"Content-Type, X-CSRF-Token")
			w.Header().Set("Access-Control-Max-Age", "600")
			w.Header().Set("Access-Control-Allow-Methods",
				"GET, POST, OPTIONS, HEAD, PUT")
		}
		next.ServeHTTP(w, r)
	}, next.Settings())
}
