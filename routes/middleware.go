package routes

import (
	"net/http"
)

// sets default content-type
func CORSMiddleware(next http.Handler, _route route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if _, allowed := allowedOrigins[origin]; allowed && _route.corsAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Max-Age", "600")
			w.Header().Set("Access-Control-Allow-Methods",
				"GET, POST, OPTIONS, HEAD, PUT")
		}
		next.ServeHTTP(w, r)
	})
}

func MethodMiddleware(next http.Handler, _route route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, found := _route.methods[r.Method]; found {
			next.ServeHTTP(w, r)
		} else if _route.corsAllowed && r.Method == "OPTIONS" {
			CORSMiddleware(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					if _, allowed := allowedOrigins[r.Header.Get("Origin")]; !allowed {
						Handle405(w, r)
					} else {
						w.WriteHeader(http.StatusOK)
					}
				}), _route).ServeHTTP(w, r)
		} else {
			Handle405(w, r)
		}
	})
}

func AuthMiddleware(next http.Handler, _route route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
