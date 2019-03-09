package routes

import (
	"database/sql"
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/forms"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/settings"
	"net/http"
)

type SigninHandler struct{}

func (h *SigninHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var form forms.SigninForm
	err := decoder.Decode(&form)
	if err != nil {
		handlers.Handle400(w, r, formats.ErrInvalidJSON, err.Error())
		return
	}

	report := form.Validate()
	if !report.Ok {
		handlers.HandleInvalidData(w, r, report, formats.ErrValidation,
			"forms.SigninForm.Validate")
		return
	}

	u, err := db.GetUserByUsernameOrEmail(settings.DB(),
		form.Username.String,
		form.Email.String)
	if err != nil {
		if err == sql.ErrNoRows {
			handlers.HandleInvalidData(w, r, forms.UnsuccessfulSigninReport,
				formats.ErrInvalidCredentials, "forms.UnsuccessfulSigninReport")
		} else {
			handlers.Handle500(w, r, formats.ErrSqlFailure,
				"db.GetUserByUsernameOrEmail", err)
		}
		return
	}

	ok, err := db.AccountComparePasswordToHash(form.Password, u.Account.Password)
	if err != nil {
		handlers.Handle500(w, r, formats.ErrPasswordHashing,
			"db.AccountComparePasswordToHash", err)
		return
	}

	if !ok {
		handlers.HandleInvalidData(w, r, forms.UnsuccessfulSigninReport,
			formats.ErrInvalidCredentials, "forms.UnsuccessfulSigninReport")
		return
	}

	err = handlers.Authorize(w, r, u)
	if err != nil {
		handlers.Handle500(w, r, formats.ErrJWTEncryptionFailure,
			"handlers.Authorize", err)
		return
	}

	handlers.Handle200(w, r, u, "db.GetUserByUsernameOrEmail")
}

func (h *SigninHandler) Methods() map[string]struct{} {
	return map[string]struct{}{"POST": {}}
}

func (h *SigninHandler) AuthRequired(method string) bool {
	return false
}

func (h *SigninHandler) CorsAllowed(method string) bool {
	return true
}