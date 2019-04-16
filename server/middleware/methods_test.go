package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/routes"
)

func TestMethods405(t *testing.T) {
	ok := false
	handlers := []routes.Handler{
		HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				t.Errorf("Request with worng method forwarded")
			},
			map[string]routes.RouteSettings{
				"GET": {
					false, true, false,
				},
			},
		),
		HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ok = true
			},
			map[string]routes.RouteSettings{
				"GET": {
					false, true, false,
				},
			},
		),
	}

	r1, _ := http.NewRequest(
		"POST", "http://localhost/Pepega", bytes.NewReader(nil))
	r2, _ := http.NewRequest(
		"GET", "http://localhost/Pepega", bytes.NewReader(nil))
	requests := []*http.Request{
		r1,
		r2,
	}

	for i, h := range handlers {
		w := httptest.NewRecorder()
		Methods(h).ServeHTTP(w, requests[i])
	}

	if !ok {
		t.Errorf("one of the handlers wasn't called")
	}
}

func TestMethodsOptions(t *testing.T) {
	handlers := []routes.Handler{
		HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
			},
			map[string]routes.RouteSettings{
				"GET": {
					false, true, false,
				},
			},
		),
		HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
			},
			map[string]routes.RouteSettings{
				"GET": {
					false, true, false,
				},
			},
		),
	}

	r1, _ := http.NewRequest(
		"OPTIONS", "http://localhost/Pepega", bytes.NewReader(nil))
	r1.Header.Set("Access-Control-Request-Method", "POST")
	r2, _ := http.NewRequest(
		"OPTIONS", "http://localhost/Pepega", bytes.NewReader(nil))
	r2.Header.Set("Access-Control-Request-Method", "GET")
	requests := []*http.Request{
		r1,
		r2,
	}

	for i, h := range handlers {
		w := httptest.NewRecorder()
		Methods(h).ServeHTTP(w, requests[i])
		if (w.Code == 405) ==
			h.Settings()[requests[i].Header.Get("Access-Control-Request-Method")].CorsAllowed {
			t.Errorf("OPTIONS request was not handled")
		}
	}
}
