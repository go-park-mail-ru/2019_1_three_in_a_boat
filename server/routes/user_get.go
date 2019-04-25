package routes

import (
	"database/sql"
	"net/http"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/formats"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/settings/shared"
)

// A handler that handles a ~single~ user resource. Simply returns the user
// data for the requested user. If the authorized user is checking his own
// profile, does not make a DB request, uses the JWE data instead. Expects GET
// method only. In case of a successful request, returns User object which is
// encoded as UserData in JSON.
func GetUser(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	uid, err := getOnlyIntParam(r.URL)
	if err != nil || uid == 0 {
		Handle404(w, r)
		return
	}

	// authedAs, ok := formats.AuthFromContext(r.Context())
	// if authedAs != nil && ok && authedAs.Pk == uid {
	// 	// if the user looking at their own profile this can save as 1 DB request
	// 	Handle200(w, r, authedAs.UserData)
	// 	return
	// }

	u, err := db.GetUser(settings.DB(), uid)
	if err != nil {
		if err == sql.ErrNoRows {
			Handle404(w, r)
		} else {
			Handle500(w, r, formats.ErrSqlFailure, err)
		}
		return
	}

	Handle200(w, r, u)
}
