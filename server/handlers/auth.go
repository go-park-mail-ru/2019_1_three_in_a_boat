package handlers

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/settings/server"
)

// Generates a JWE authorization token for user. Sets the cookie. Returns error
// if the token generation failed (shouldn't ever happen in a properly
// configured app.
func Authorize(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * server_settings.JWETokenLifespan),
		HttpOnly: true,
		Path:     "/", // guarantees uniqueness.. I think
	})
}

// Deletes the authorization cookie. Guaranteed to succeed.
func Unauthorize(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Expires:  time.Now().Add(-24 * time.Hour), // sufficient for any time zone
		HttpOnly: true,
		Path:     "/", // guarantees uniqueness.. I think
	})
}
