package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/shared"
)

func TestCSRF(t *testing.T) {
	//noinspection GoBoolExpressions
	if !settings.EnableCSRF {
		fmt.Println("CSRF is not enabled: skipping csrf test")
		return
	}
	ok := false
	handlers := []http_utils.Handler{
		HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				t.Errorf("Request with bad CSRF forwarded")
			},
			map[string]http_utils.RouteSettings{
				"GET": {
					false, true, true,
				},
			},
		),
		HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ok = true
			},
			map[string]http_utils.RouteSettings{
				"GET": {
					false, true, true,
				},
			},
		),
	}

	r1, _ := http.NewRequest(
		"GET", "http://localhost/Pepega", bytes.NewReader(nil))
	r1.AddCookie(&http.Cookie{Name: "csrf", Value: "foobar"})
	r1.Header.Set("X-CSRF-Token", "notfoobar")
	r2, _ := http.NewRequest(
		"GET", "http://localhost/Pepega", bytes.NewReader(nil))
	r2.AddCookie(&http.Cookie{Name: "csrf", Value: "foobar"})
	r2.Header.Set("X-CSRF-Token", "foobar")
	requests := []*http.Request{
		r1,
		r2,
	}

	for i, h := range handlers {
		w := httptest.NewRecorder()
		CSRF(h).ServeHTTP(w, requests[i])
	}

	if !ok {
		t.Errorf("one of the handlers wasn't called")
	}
}
