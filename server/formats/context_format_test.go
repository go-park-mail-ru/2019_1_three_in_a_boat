package formats

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/db"
)

func TestAuthContext(t *testing.T) {
	ctx := context.Background()
	uc := UserClaims{
		UserData: &db.UserData{
			Pk: 42,
		},
	}
	ctx = NewAuthContext(ctx, &uc)

	uc2, ok := AuthFromContext(ctx)

	if !ok {
		t.Error("failed to extract auth from context after adding it")
	} else if uc2.Pk != 42 {
		t.Error("extracted empty/corrupted userdata")
	}

	_, ok = AuthFromContext(context.Background())
	if ok {
		t.Error("extracted auth from empty context")
	}
}
