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

func Authorize(w http.ResponseWriter, r *http.Request, user *db.User) error {
	token, err := makeJWEToken(user)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
	})

	return nil
}

func makeJWEToken(user *db.User) (string, error) {
	claimUid, err := uuid.NewRandom()
	if err != nil {
		return "", err
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
			Img:        user.Profile.Img,
			BirthDate:  user.Profile.BirthDate,
			SignupDate: user.Profile.SignupDate,
		},
	}

	builder := jwt.Signed(settings.GetSigner())
	builder = builder.Claims(customClaims)

	return builder.CompactSerialize()

}
