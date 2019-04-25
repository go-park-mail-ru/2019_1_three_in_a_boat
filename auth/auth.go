package main

import (
	"github.com/google/uuid"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/formats/pb"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/settings/auth"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/settings/shared"
)

// Generates a JWE token for user that can be saved in the cookies.
func tokenizeUser(user *db.User) (string, error) {
	return tokenizeClaims(formats.ClaimsFromUser(user))
}

func tokenizeClaims(claims *pb.Claims) (string, error) {
	claimUid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	claims.TokenId = claimUid.String()
	if claims.Img == "" {
		claims.Img = settings.DefaultImgName
	}

	builder := jwt.Signed(auth_settings.GetSigner())
	builder = builder.Claims(claims)
	return builder.CompactSerialize()
}
