package formats

import (
	"context"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats/pb"
)

// The file provides interface to safely interact with Context

type key int

var userKey key

// Adds given *UserClaims to the given context using a guaranteed unique key
func NewAuthContext(ctx context.Context, u *pb.Claims) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func ClaimsFromUser(u *db.User) *pb.Claims {
	return &pb.Claims{
		Uid:      u.Account.Pk,
		Username: u.Account.Username,
		Email:    u.Account.Email,
		Score:    u.Profile.HighScore.Int64,
		Gender:   u.Profile.Gender.String,
		Img:      u.Profile.Img.String,
	}
}

// Retrieves UserClaims from the context using the same key as NewAuthContext.
// Second boolean value indicates whether the type assertion succeeded, which
// should always happen. UserClaims can still be nil if bool is true.
func AuthFromContext(ctx context.Context) (*pb.Claims, bool) {
	u, ok := ctx.Value(userKey).(*pb.Claims)
	return u, ok
}
