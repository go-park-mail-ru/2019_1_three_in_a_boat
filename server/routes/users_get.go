package routes

import (
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
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
	if HandleErrForward(w, r, formats.ErrSqlFailure, err) != nil {
		return
	} else {
		defer rows.Close()
	}

	users := GetUsersResponse{}
	var nUsers int
	var u *db.User
	for i := 0; rows.Next() && i < pageSize; i++ {
		u, nUsers, err = db.UserFromRow(rows)
		if HandleErrForward(w, r, formats.ErrDbScanFailure, err) != nil {
			return
		}
		users.Users = append(users.Users, u)
	}

	if err := rows.Err(); HandleErrForward(w, r, formats.ErrDbRowsFailure, err) != nil {
		return
	}

	if nUsers == 0 {
		Handle404(w, r)
		return
	}

	users.Page = page
	users.NPages = (nUsers - 1) / pageSize

	Handle200(w, r, users)
}

func validateUsersParams(url *url.URL) (
	page int, pageSize int, order []db.SelectOrder) {
	return makeNPage(url.Query()["page"]),
		makePageSize(url.Query()["pageSize"], UsersMaxPageSize, UsersDefaultPageSize),
		makeOrderList(url.Query()["sort"], db.UserOrderMap)
}
