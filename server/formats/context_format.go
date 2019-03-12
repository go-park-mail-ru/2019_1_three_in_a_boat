package formats

import (
	"context"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/db"
	"gopkg.in/square/go-jose.v2/jwt"
)

// file provides interface to safely interact with Context

type key int

var userKey key

// Represents JWT claims, as stored in the JWT and in the context
type UserClaims struct {
	*jwt.Claims
	*db.UserData
}

// Adds given *UserClaims to the given context using a guaranteed unique key
func NewAuthContext(ctx context.Context, u *UserClaims) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// Retrieves UserClaims from the context using the same key as NewAuthContext.
// Second boolean value indicates whether the type assertion succeeded, which
// should always happen. UserClaims can still be nil if bool is true.
func AuthFromContext(ctx context.Context) (*UserClaims, bool) {
	u, ok := ctx.Value(userKey).(*UserClaims)
	return u, ok
}
