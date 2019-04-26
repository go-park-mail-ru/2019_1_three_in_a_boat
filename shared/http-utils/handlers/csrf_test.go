package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSetNewCsrfToken(t *testing.T) {
	w := httptest.NewRecorder()

	err := SetNewCsrfToken(w)
	if err != nil {
		t.Fatal("failed to write token")
	}

	var cookie *http.Cookie = nil
	for _, c := range w.Result().Cookies() {
		if c.Name == "csrf" {
			cookie = c
			break
		}
	}

	if cookie == nil {
		t.Fatal("CSRF cookie was not set")
	}

	if cookie.HttpOnly {
		t.Error("CSRF cookie is HTTP-only and can not be read by a script")
	}

	if cookie.Expires.Before(time.Now()) {
		t.Error("CSRF cookie expires too quickly")
	}

}
