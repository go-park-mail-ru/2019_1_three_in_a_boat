package routes

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/forms"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/settings"
	"net/http"
)

// stores pointers, assumes they aren't modified anywhere else
type PutUserResponse = db.User

// A handler that handles a ~single~ user resource. Requires authorization,
// returns 403 if the user is not authorized or is trying to change somebody
// else's data. Uses UserEditForm to validate the data. Handles saving and
// deleting images, however, the resizing and ID generation is handled in the
// form. Expects PUT method only. In case of a successful request, returns
// UserData with the updated data. The request.user is going to be updated
// on the following request. In case of a failure, returns form.Report
// indicating errors in the user data.
func PutUser(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	uid, err := getOnlyIntParam(r.URL)
	if err != nil || uid == 0 {
		Handle404(w, r)
		return
	}

	authedAs, ok := formats.AuthFromContext(r.Context())
	// since AuthRequired is true, authedAs is always a valid UserClaims
	// let's just check though LUL. Also authedAs.Pk != uid is a necessary check.
	if authedAs == nil || !ok || authedAs.Pk != uid {
		Handle403(w, r)
		return
	}

	form := forms.UserEditForm{}
	err = json.NewDecoder(r.Body).Decode(&form)
	if HandleErrForward(w, r, formats.ErrInvalidJSON, err) != nil {
		return
	}

	report := form.Validate()
	if !HandleReportForward(w, r, report).Ok {
		return
	}

	// might not be needed, but is at least convenient
	u, err := db.GetUser(settings.DB(), authedAs.Pk)
	if HandleErrForward(w, r, formats.ErrSqlFailure, err) != nil {
		return
	}

	oldImg := u.Profile.Img
	u, err = form.EditUser(u)
	if HandleErrForward(w, r, formats.ErrDbModuleFailure, err) != nil {
		return
	}

	err, txErr := u.Save(settings.DB())
	if HandleErrForward(w, r, formats.ErrSavingImg, txErr) != nil {
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

	if form.Img != nil {
		// here's a little bug that can't be easily fixed: we want to create a file
		// and write to the database transactionally, which isn't really possible.
		// since writing to a file is extremely unlikely to fail (no disk space or
		// something of the sort), we do it second. If this fails, we just return an
		// error, and the user will be left with a broken link instead of a pic.
		// This, however, can be easily fixed by uploading the same pic again,
		// assuming the error will go away. So we just let it be.
		err = SaveImage(form.Img, u.Profile.Img.String)
		if HandleErrForward(w, r, formats.ErrSavingImg, err) != nil {
			return
		}
	}
	if form.Img != nil && oldImg.Valid {
		// we don't care if this fails and this is very unlikely to fail
		_ = DeleteImage(oldImg.String)
	}

	if HandleErrForward(w, r, formats.ErrSignupAuthFailure, Authorize(w, u)) != nil {
		return
	}

	Handle200(w, r, u)
}
