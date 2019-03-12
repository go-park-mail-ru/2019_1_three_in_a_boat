package handlers

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/settings"
	"github.com/google/uuid"
	"gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"time"
)

// Generates a JWE authorization token for user. Sets the cookie. Returns error
// if the token generation failed (shouldn't ever happen in a properly
// configured app.
func Authorize(w http.ResponseWriter, user *db.User) error {
	token, err := makeJWEToken(user)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * settings.JWETokenLifespan),
		HttpOnly: true,
		Path:     "/", // guarantees uniqueness.. I think
	})

	return nil
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

// Generates a JWE token for user that can be saved in the cookies.
func makeJWEToken(user *db.User) (string, error) {
	claimUid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	img := user.Profile.Img
	if !img.Valid || img.String == "" {
		img.String = "default.png"
		img.Valid = true
	}

	customClaims := formats.UserClaims{
		Claims: &jwt.Claims{
			Issuer:   "hexagon",
			Subject:  user.Account.Username,
			ID:       claimUid.String(),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
		UserData: &db.UserData{
			Pk:         user.Account.Pk,
			Username:   user.Account.Username,
			Email:      user.Account.Email,
			FirstName:  user.Profile.FirstName,
			LastName:   user.Profile.LastName,
			HighScore:  user.Profile.HighScore,
			Gender:     user.Profile.Gender,
			Img:        img,
			BirthDate:  db.NullDateTime{user.Profile.BirthDate},
			SignupDate: user.Profile.SignupDate,
		},
	}

	builder := jwt.Signed(settings.GetSigner())
	builder = builder.Claims(customClaims)

	return builder.CompactSerialize()
}
