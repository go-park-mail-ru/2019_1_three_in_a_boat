package routes

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/forms"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/settings"
	"net/http"
)

type PostUsersResponse = forms.SignupForm

// Handler for the Users resource
func PostUsers(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var form forms.SignupForm
	err := decoder.Decode(&form)
	if err != nil {
		handlers.Handle400(w, r, formats.ErrInvalidJSON, err.Error())
		return
	}

	report := form.Validate()
	if !report.Ok {
		handlers.HandleInvalidData(w, r, report, formats.ErrValidation,
			"forms.SignupForm.Validate")
		return
	}

	u, err := form.MakeUser()
	if err != nil {
		handlers.Handle500(w, r, formats.ErrDbModuleFailure,
			"forms.SignupForm.MakeUser", err)
		return
	}

	err, txErr := u.Save(settings.DB())
	if txErr != nil {
		handlers.Handle500(w, r, formats.ErrDbTransactionFailure,
			"db.User.Save", txErr)
		return
	}

	report, err = forms.CheckUserDbConstraints(err)
	if err != nil {
		if report != nil {
			handlers.HandleInvalidData(w, r, report,
				formats.ErrValidation, "forms.NewReport")
		} else {
			handlers.Handle500(w, r, formats.ErrSqlFailure, "db.User.Save", err)
		}
		return
	}

	err = handlers.Authorize(w, r, u)
	if err != nil {
		handlers.Handle500(w, r, formats.ErrSignupAuthFailure,
			"handlers.Authorize", err)
		return
	}

	handlers.Handle201(w, r, u, "forms.SignupForm.MakeUser")
}
