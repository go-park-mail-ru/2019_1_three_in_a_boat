package routes

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/test-utils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type GetUserFullResponse struct {
	formats.JSONResponseData
	Data db.UserData `json:"data"`
}

func TestGetUser(t *testing.T) {
	SetUpDB()
	defer TearDownDB()
	cases := []struct {
		req  func(m *Mock) *http.Request
		res  func(m *Mock) db.UserData
		fail bool
	}{
		{
			req: func(m *Mock) *http.Request {
				return httptest.NewRequest(
					"GET",
					fmt.Sprintf("http://localhost/users/%d", m.GetPk()),
					nil)
			},
			res: func(m *Mock) db.UserData {
				return MockToUserData(m)
			},
			fail: false,
		},
	}

	for _, c := range cases {
		for _, mockUser := range GetMockData() {
			w := httptest.NewRecorder()
			GetUser(w, c.req(&mockUser))
			body, _ := ioutil.ReadAll(w.Body)
			var res GetUserFullResponse
			if err := json.Unmarshal(body, &res); err != nil {
				t.Error(err)
			}

			if formats.StatusMap[!c.fail] != res.Status {
				t.Errorf("status mismatch: %s != %s",
					formats.StatusMap[!c.fail], res.Status)
			}

			ud := c.res(&mockUser)
			if !UserDataEqual(res.Data, ud) {
				t.Errorf("response mismatch: %v != %v", res.Data, ud)
			}
		}
	}
}
