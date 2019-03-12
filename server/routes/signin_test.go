package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/test-utils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type SigninResponse struct {
	Data db.UserData `json:"data"`
}

type FullSigninResponse struct {
	formats.JSONResponseData
	Data db.UserData `json:"data"`
}

type SignInForm struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func TestSigninHandler_ServeHTTP(t *testing.T) {
	cases := []struct {
		req  func() *http.Request
		res  func() SigninResponse
		fail bool
	}{
		{
			req: func() *http.Request {
				form, _ := json.Marshal(SignInForm{
					GetMockData()[0].Account.Username,
					"12345"})
				return httptest.NewRequest(
					"GET",
					"http://localhost/signin",
					bytes.NewReader(form))
			},
			res: func() SigninResponse {
				return SigninResponse{Data: UserToUserData(GetMockData()[0].User)}
			},
			fail: false,
		},
		{
			req: func() *http.Request {
				form, _ := json.Marshal(SignInForm{
					GetMockData()[0].Account.Username,
					"123456"})
				return httptest.NewRequest(
					"GET",
					"http://localhost/signin",
					bytes.NewReader(form))
			},
			res: func() SigninResponse {
				return SigninResponse{}
			},
			fail: true,
		},
		{
			req: func() *http.Request {
				form, _ := json.Marshal(SignInForm{
					GetMockData()[0].Account.Email,
					"12345"})
				return httptest.NewRequest(
					"GET",
					"http://localhost/signin",
					bytes.NewReader(form))
			},
			res: func() SigninResponse {
				return SigninResponse{Data: UserToUserData(GetMockData()[0].User)}
			},
			fail: false,
		},
	}

	for _, c := range cases {
		w := httptest.NewRecorder()
		handler := SigninHandler{}
		handler.ServeHTTP(w, c.req())
		body, _ := ioutil.ReadAll(w.Body)
		var res FullSigninResponse
		if err := json.Unmarshal(body, &res); err != nil {
			t.Error(err)
		}

		if formats.StatusMap[!c.fail] != res.Status {
			t.Errorf("status mismatch: %s != %s",
				formats.StatusMap[!c.fail], res.Status)
		}

		authCookieFound := false
		for _, cookie := range w.Result().Cookies() {
			if cookie.Name == "auth" && cookie.Value != "" {
				authCookieFound = true
				break
			}
		}

		if authCookieFound == c.fail {
			if c.fail {
				t.Error("Unexpected auth cookie")
			} else {
				t.Error("Auth cookie not found")
			}
		}

		if !c.fail && res.Data.Pk == 0 {
			t.Error("Auth was not supposed to fail but dir")
		} else if c.fail && res.Data.Pk != 0 {
			t.Error("Auth was supposed to fail but returned a user")
		}

	}
}
