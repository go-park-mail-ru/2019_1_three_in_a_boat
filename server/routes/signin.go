package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/forms"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/formats/pb"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/http-utils/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/shared"
)

// Handles Signin resource. Only accepts POST requests. Implements
// routes.Handler interface, which extends http.Handler. Uses SigninForm to
// validate the data. In case of a successful response, returns User which gets
// encoded into JSON as db.UserData
type SigninHandler struct{}

// Handles POST requests.
func (h *SigninHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var form forms.SigninForm
	err := decoder.Decode(&form)
	if err != nil {
		Handle400(w, r, formats.ErrInvalidJSON, err.Error())
		return
	}

	report := form.Validate()
	if !HandleReportForward(w, r, report).Ok {
		return
	}

	if ok, claims := authorize(
		w, r, form.Username.String, form.Email.String, form.Password); ok {
		Handle200(w, r, claims)
	} // else do nothing - authorize handles errors itself

}

func (h *SigninHandler) Settings() map[string]http_utils.RouteSettings {
	return map[string]http_utils.RouteSettings{
		"POST": {
			AuthRequired:           false,
			CorsAllowed:            true,
			CsrfProtectionRequired: true,
		},
	}
}

func authorize(
	w http.ResponseWriter, r *http.Request, username, email, password string) (
	bool, *pb.Claims) {
	authCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	checkAuthReply, err := settings.AuthClient.Authorize(authCtx,
		&pb.AuthorizeRequest{
			Username: username,
			Email:    email,
			Password: password,
		})

	if err != nil {
		if checkAuthReply != nil {
			Handle500(w, r, checkAuthReply.Message, err)
		} else {
			Handle500(w, r, formats.ErrAuthServiceFailure, err)
		}
		return false, nil
	}

	if !checkAuthReply.GetOk() {
		HandleInvalidData(w, r, forms.UnsuccessfulSigninReport, formats.ErrInvalidCredentials)
		return false, nil
	}

	// if no errors, message is the token
	Authorize(w, checkAuthReply.Message)
	return true, checkAuthReply.Claims
}
