package middleware

import (
	"context"
	"net/http"
	"time"

	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/routes"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/settings"
)

// Authentication middleware: if the resource requires authentication,
// checks JWT token and calls the method, returning 403 if it's invalid.
// Otherwise simply forwards the call. MUST be used after every other middleware,
// since it unconditionally adds data to the request. That data is public,
// however, you still wouldn't want to give it away in a CSRF attack.  ALWAYS
// adds formats.UserClaims to the context - default constructed if no user.
// Verifies that the claims are signed.
func Auth(next routes.Handler) routes.Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		var rawJWT *http.Cookie
		var parsedJWT *jwt.JSONWebToken
		var err error
		errMsg := ""

		rawJWT, err = r.Cookie("auth")
		if err != nil { // err = ErrNoCookie only
			ctx = formats.NewAuthContext(context.Background(), nil)
		} else {
			parsedJWT, err = jwt.ParseSigned(rawJWT.Value)
			if err != nil {
				ctx = formats.NewAuthContext(context.Background(), nil)
				errMsg = formats.ErrJWTDecryptionFailure
			} else {
				claims := formats.UserClaims{}
				err = parsedJWT.Claims(&settings.GetSecretKey().PublicKey, &claims)
				if err != nil {
					errMsg = formats.ErrJWTDecryptionFailure
				} else if claims.Pk == 0 {
					errMsg = formats.ErrJWTDecryptionEmpty
				} else if claims.Expiry.Time().Before(time.Now()) {
					errMsg = formats.ErrJWTOutdated
				} else {
					ctx = formats.NewAuthContext(context.Background(), &claims)
				}
			}
		}

		// err indicates any failure in the above functions, while errMsg is only
		// written if something went unexpectedly wrong (i.e., jwt key present, but
		// could not be read)
		if errMsg != "" {
			ctx = formats.NewAuthContext(context.Background(), nil)
			//noinspection GoNilness
			handlers.LogError(0, errMsg+": "+err.Error(), r)
			// log the error and forward the response as unauthorized - do not return
		}

		if next.Settings()[r.Method].AuthRequired && err != nil {
			handlers.LogError(0, err.Error(), r)
			handlers.Handle403(w, r)
			return // .. unless the resource requires authorization
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}, next.Settings())
}
