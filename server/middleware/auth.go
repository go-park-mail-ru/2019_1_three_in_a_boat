package middleware

import (
	"context"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/routes"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/settings"
	"gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"time"
)

// Authentication middleware: if the resource requires authentication,
// checks JWT token and calls the method, returning 403 if it's invalid.
// Otherwise simply forwards the call.
// MUST be used after Methods middleware, or it might silently forward
// unauthorized requests.
// ALWAYS adds formats.UserClaims to the context - default constructed if no
// user. Checks the claims
func Auth(next http.Handler, _route routes.Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			handlers.LogError(0, errMsg, r)
			// log the error and forward the response as unauthorized - do not return
		}
		if _route.Handler.AuthRequired(r.Method) {
			handlers.Handle403(w, r)
			return // .. unless the resource requires authorization
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
		//if !_route.Handler.AuthRequired(r.Method) {
		//	r.WithContext(ctx)
		//
		//} else {
		//
		//}
	})
}
