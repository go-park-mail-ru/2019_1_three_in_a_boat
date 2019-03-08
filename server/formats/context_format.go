package formats

import (
	"context"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"gopkg.in/square/go-jose.v2/jwt"
)

type key int

var userKey key

type UserClaims struct {
	*jwt.Claims
	*db.UserData
}

func NewAuthContext(ctx context.Context, u *UserClaims) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func AuthFromContext(ctx context.Context) (*UserClaims, bool) {
	u, ok := ctx.Value(userKey).(*UserClaims)
	return u, ok
}
