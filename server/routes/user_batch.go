package routes

import (
	"errors"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/shared"
	"net/http"
	"net/url"
	"strconv"
)

// Represents the result returned by GetUsers
type GetUsersBatchResponse struct {
	Users []*db.User `json:"users"` // will be transformed to UserData
}

type UsersBatchHandler struct{}

// A handler that handles a ~multiple~ users resource. Only accepts GET
// requests. Accepts sort=[]string, page=int, pageSize=int GET params validates
// those and sends them to DB. Unlike forms, incorrect parameters are ignored
// rather than triggering an error. In case of a successful request, returns
// GetUsersResponse struct, described above.
func (*UsersBatchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	ids, err := getIdsFromUri(r.URL)
	if err != nil {
		HandleInvalidData(w, r, formats.ErrInvalidGetParams, err.Error())
		return
	}

	rows, err := db.GetBatchByIds(settings.DB(), ids)
	if HandleErrForward(w, r, formats.ErrSqlFailure, err) != nil {
		return
	}

	//noinspection GoUnhandledErrorResult
	defer rows.Close()
	users := GetUsersBatchResponse{}
	for rows.Next() {
		u, err := db.UserFromBatchRow(rows)
		if HandleErrForward(w, r, formats.ErrDbScanFailure, err) != nil {
			return
		}
		users.Users = append(users.Users, u)
	}

	if err := rows.Err(); HandleErrForward(w, r, formats.ErrDbRowsFailure, err) != nil {
		return
	}

	Handle200(w, r, users)

}

func (h *UsersBatchHandler) Settings() map[string]http_utils.RouteSettings {
	return map[string]http_utils.RouteSettings{
		"GET": {
			AuthRequired:           false,
			CorsAllowed:            true,
			CsrfProtectionRequired: false,
		},
	}
}

func getIdsFromUri(url *url.URL) ([]int64, error) {
	if len(url.Query()["ids"]) == 0 {
		return nil, errors.New("no ids")
	}
	var res []int64
	for _, id := range url.Query()["ids"] {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		res = append(res, int64(idInt))
	}

	return res, nil
}
