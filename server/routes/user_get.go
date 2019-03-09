package routes

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/settings"
	"net/http"
)

// stores pointers, assumes they aren't modified anywhere else
type GetUserResponse = db.User

// Handler for the Users resource
func GetUser(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	uid, err := getOnlyIntParam(r.URL)
	if err != nil || uid == 0 {
		Handle404(w, r)
		return
	}

	authedAs, ok := formats.AuthFromContext(r.Context())
	if authedAs != nil && ok && authedAs.Pk == uid {
		// if the user looking at their own profile this can save as 1 DB request
		Handle200(w, r, authedAs.UserData)
		return
	}

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
