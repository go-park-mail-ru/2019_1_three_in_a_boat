package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/shared"
)

func SetNewCsrfToken(w http.ResponseWriter) error {
	token, err := MakeCSRFToken()
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * settings.CSRFTokenLifespan),
		HttpOnly: false,
		Path:     "/", // guarantees uniqueness.. I think
	})

	return nil
}

func MakeCSRFToken() (string, error) {
	randombytes := make([]byte, settings.CSRFTokenLength)
	_, err := rand.Read(randombytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(randombytes), nil
}
