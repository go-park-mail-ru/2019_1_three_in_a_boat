package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/server"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/test-utils"
)

type GetUsersFullResponse struct {
	formats.JSONResponseData
	Data GetUsersResponseUnmarhsaled `json:"data"`
}
type GetUsersResponseUnmarhsaled struct {
	Users  []db.UserData `json:"users"`  // will be transformed to UserData
	Page   int           `json:"page"`   // 0-indexed
	NPages int           `json:"nPages"` // largest valid Page value is NPages - 1
}

func TestGetUsers(t *testing.T) {
	cases := []struct {
		req  *http.Request
		res  func() GetUsersResponse
		fail bool
	}{
		{
			req: httptest.NewRequest(
				"GET",
				"http://localhost/users?sort=-pk",
				nil),
			res: func() GetUsersResponse {
				sorted := make([]Mock, len(GetMockData()))
				copy(sorted, GetMockData())
				sort.Slice(sorted, func(i, j int) bool {
					return sorted[i].GetPk() > sorted[j].GetPk()
				})

				res := make([]*db.User, len(GetMockData()))
				for i := range sorted {
					res[i] = sorted[i].User
				}

				return GetUsersResponse{
					res,
					0,
					(len(sorted) - 1) / server_settings.UsersDefaultPageSize,
				}
			},
			fail: false,
		},
		{
			req: httptest.NewRequest(
				"GET",
				"http://localhost/users?sort=pk&page=0&pageSize=2",
				nil),
			res: func() GetUsersResponse {
				sorted := make([]Mock, len(GetMockData()))
				copy(sorted, GetMockData())
				sort.Slice(sorted, func(i, j int) bool {
					return sorted[i].GetPk() < sorted[j].GetPk()
				})
				res := make([]*db.User, len(sorted[:2]))
				for i := range sorted[:2] {
					res[i] = sorted[:2][i].User
				}

				return GetUsersResponse{
					res,
					0,
					(len(sorted) - 1) / 2,
				}

			},
			fail: false,
		},
		{
			req: httptest.NewRequest(
				"GET",
				"http://localhost/users?sort=pk&page=1&pageSize=2",
				nil),
			res: func() GetUsersResponse {
				sorted := make([]Mock, len(GetMockData()))
				copy(sorted, GetMockData())
				sort.Slice(sorted, func(i, j int) bool {
					return sorted[i].GetPk() < sorted[j].GetPk()
				})

				res := make([]*db.User, len(sorted[2:]))
				for i := range sorted[2:] {
					res[i] = sorted[2:][i].User
				}

				return GetUsersResponse{
					res,
					1,
					(len(sorted) - 1) / 2,
				}
			},
			fail: false,
		},
		{
			req: httptest.NewRequest(
				"GET",
				"http://localhost/users?sort=-username&pageSize=2",
				nil),
			res: func() GetUsersResponse {
				sorted := make([]Mock, len(GetMockData()))
				copy(sorted, GetMockData())
				sort.Slice(sorted, func(i, j int) bool {
					return strings.ToLower(sorted[i].Account.Username) >
						strings.ToLower(sorted[j].Account.Username)
				})
				res := make([]*db.User, len(sorted[:2]))
				for i := range sorted[:2] {
					res[i] = sorted[:2][i].User
				}

				return GetUsersResponse{
					res,
					0,
					(len(sorted) - 1) / 2,
				}
			},
			fail: false,
		},
		{
			req: httptest.NewRequest(
				"GET",
				"http://localhost/users?sort=-username&pageSize=10&page=2",
				nil),
			res:  func() GetUsersResponse { return GetUsersResponse{} },
			fail: true,
		},
	}

	for i, c := range cases {
		w := httptest.NewRecorder()
		GetUsers(w, c.req)
		body, _ := ioutil.ReadAll(w.Body)
		var res GetUsersFullResponse
		expected := c.res()
		if err := json.Unmarshal(body, &res); err != nil {
			t.Error(err)
		}

		if formats.StatusMap[!c.fail] != res.Status {
			t.Errorf("case %d: unexpected failure or success: %s != %s", i,
				formats.StatusMap[c.fail], res.Status)
		}

		if expected.Page != res.Data.Page || expected.NPages != res.Data.NPages {
			t.Errorf("case %d: pagination mismatch: %d, %d != %d, %d", i,
				expected.Page, expected.NPages, res.Data.Page, res.Data.NPages)
		}

		if len(res.Data.Users) != len(expected.Users) {
			t.Errorf("Array len differs: %d != %d",
				len(res.Data.Users), len(expected.Users))
		} else {
			for j := range expected.Users {
				mockUD := UserToUserData(expected.Users[j])
				resUD := res.Data.Users[j]
				if !UserDataEqual(resUD, mockUD) {
					t.Errorf("case %d: entry %d differs: %v != %v", i, j, resUD, mockUD)
				}
			}
		}
	}
}
