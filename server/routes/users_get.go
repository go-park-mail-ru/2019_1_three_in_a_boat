package routes

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/settings"
	"net/http"
	"net/url"
)

// stores pointers, assumes they aren't modified anywhere else
type GetUsersResponse struct {
	Users  []*db.User `json:"users"`
	Page   int        `json:"page"`   // 0-indexed
	NPages int        `json:"nPages"` // largest valid Page value is NPages - 1
}

// TODO: split this in a separate routes/settings file if more constants arise
const UsersDefaultPageSize = 10
const UsersMaxPageSize = 10

// Handler for the Users resource
func GetUsers(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	page, pageSize, order := validateUsersParams(r.URL)

	rows, err := db.GetUserMany(settings.DB(), order, -1, pageSize*page)
	if err != nil {
		handlers.Handle500(w, r, formats.ErrSqlFailure, "db.GetUserMany", err)
		return
	} else {
		defer rows.Close()
	}

	users := GetUsersResponse{}
	var nUsers int
	var u *db.User
	for i := 0; rows.Next() && i < pageSize; i++ {
		u, nUsers, err = db.UserFromRow(rows)
		if err != nil {
			handlers.Handle500(w, r, formats.ErrDbScanFailure,
				"db.UserFromRow", err)
			return
		}
		users.Users = append(users.Users, u)
	}

	if err := rows.Err(); err != nil {
		handlers.Handle500(w, r, formats.ErrDbRowsFailure, "db.GetUserMany", err)
		return
	}

	if nUsers == 0 {
		handlers.Handle404(w, r)
		return
	}

	users.Page = page
	users.NPages = (nUsers - 1) / pageSize

	handlers.Handle200(w, r, users, "db.AuthorDataFromRow")
}

func validateUsersParams(url *url.URL) (
	page int, pageSize int, order []db.SelectOrder) {
	return makeNPage(url.Query()["page"]),
		makePageSize(url.Query()["pageSize"], UsersMaxPageSize, UsersDefaultPageSize),
		makeOrderList(url.Query()["sort"], db.UserOrderMap)
}
