package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats/pb"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/shared"
)

// Authentication middleware: if the resource requires authentication,
// checks JWT token and calls the method, returning 403 if it's invalid.
// Otherwise simply forwards the call. MUST be used after every other middleware,
// since it unconditionally adds data to the request. That data is public,
// however, you still wouldn't want to give it away in a CSRF attack.  ALWAYS
// adds formats.UserClaims to the context - default constructed if no user.
// Verifies that the claims are signed.
func Auth(next http_utils.Handler) http_utils.Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var rawJWT *http.Cookie
		var err error

		rawJWT, err = r.Cookie("auth")
		if err != nil { // err = ErrNoCookie only
			handleAuthError(w, r, next)
			return
		}

		authCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		authReply, err := settings.AuthClient.CheckAuthorize(
			authCtx, &pb.CheckAuthorizeRequest{Token: rawJWT.Value})

		if err != nil {
			if authReply != nil {
				handlers.LogError(0, authReply.Message+": "+err.Error(), r)
			} else {
				handlers.LogError(0, err.Error(), r)
			}
			handleAuthError(w, r, next)
			return
		}

		ctx = formats.NewAuthContext(ctx, authReply.Claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}, next.Settings())
}

// Initializes an empty context, calls next if auth is not required, handles 403 otherwise
func handleAuthError(
	w http.ResponseWriter,
	r *http.Request,
	next http_utils.Handler) {
	ctx := formats.NewAuthContext(context.Background(), nil)
	r = r.WithContext(ctx)
	if next.Settings()[r.Method].AuthRequired {
		handlers.Handle403(w, r)
	} else {
		next.ServeHTTP(w, r)
	}
}
