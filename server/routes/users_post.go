package routes

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/forms"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
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
		Handle400(w, r, formats.ErrInvalidJSON, err.Error())
		return
	}

	report := form.Validate()
	if !HandleReportForward(w, r, report).Ok {
		return
	}

	u, err := form.MakeUser()
	if HandleErrForward(w, r, formats.ErrDbModuleFailure, err) != nil {
		return
	}

	err, txErr := u.Save(settings.DB())
	if HandleErrForward(w, r, formats.ErrDbTransactionFailure, txErr) != nil {
		return
	}

	report, err = forms.CheckUserDbConstraints(err)
	if report == nil {
		if HandleErrForward(w, r, formats.ErrSqlFailure, err) != nil {
			return
		}
	} else if !HandleReportForward(w, r, report).Ok {
		return
	}

	if HandleErrForward(w, r, formats.ErrSignupAuthFailure, Authorize(w, r, u)) != nil {
		return
	}

	Handle201(w, r, u)
}
