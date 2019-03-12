package routes

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/test-utils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Tests are dependent on DB because I'm a stupid singleton-using idiot.
// That being said, writing a DB Mock was harder than working with postgres so
// I'm kinda k with how it all turned out.

type AuthorsFullResponse struct {
	formats.JSONResponseData
	Data []db.AuthorData `json:"data"`
}

func TestAuthorsHandler_ServeHTTPs(t *testing.T) {
	SetUpDB()
	defer TearDownDB()
	cases := []struct {
		req *http.Request
		ok  func(*httptest.ResponseRecorder)
	}{
		{
			req: httptest.NewRequest("GET", "http://localhost/authors", nil),
			ok: func(resp *httptest.ResponseRecorder) {
				body, _ := ioutil.ReadAll(resp.Body)
				var res AuthorsFullResponse
				if err := json.Unmarshal(body, &res); err != nil {
					t.Error(err)
				}

				if res.Status != formats.StatusMap[true] {
					t.Error("Error returned by the handler!")
				}

				if len(GetMockData()) != len(res.Data) {
					t.Errorf(
						"duplicate or not all authors have been returned: %d != %d",
						len(res.Data), len(GetMockData()))
				}

				for mock := range GetMockData() {
					found := false
					for au := range res.Data {
						if au == mock {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("%v not in the result", mock)
					}
				}

			},
		},
	}

	handler := AuthorsHandler{}
	for _, c := range cases {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, c.req)
		c.ok(w)
	}
}
