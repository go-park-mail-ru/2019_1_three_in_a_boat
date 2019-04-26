package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPanic(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		panic("DEADLINE IS TOMORROW TESTS DONT WORK THE WALLS ARE CLOSING IN AAAAA")
	}

	w := httptest.NewRecorder()
	r, _ :=
		http.NewRequest("GET", "http://localhost/Pepega", bytes.NewReader(nil))
	Panic(HandlerFunc(handler, nil)).ServeHTTP(w, r) // if it panics the test fails
}
