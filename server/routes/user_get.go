package routes

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/settings"
	"net/http"
	"strconv"
	"strings"
)

// stores pointers, assumes they aren't modified anywhere else
type GetUserResponse = db.User

// Handler for the Users resource
func GetUser(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	urlParamsSplit := strings.Split(r.URL.Path[1:], "/")
	if len(urlParamsSplit) != 2 {
		handlers.Handle404(w, r)
		return
	}

	uid, err := strconv.ParseInt(urlParamsSplit[1], 10, 64)
	if err != nil || uid == 0 {
		handlers.Handle404(w, r)
		return
	}

	u, err := db.GetUser(settings.DB(), uid)
	if err != nil {
		if err == sql.ErrNoRows {
			handlers.Handle404(w, r)
		} else {
			handlers.Handle500(w, r, formats.ErrSqlFailure, "db.GetUser", err)
		}
		return
	}

	handlers.Handle200(w, r, u, "db.AuthorDataFromRow")
}
