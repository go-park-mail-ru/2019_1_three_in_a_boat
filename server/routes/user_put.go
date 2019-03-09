package routes

import (
	"encoding/json"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/forms"
	. "github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/handlers"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/settings"
	"net/http"
	"os"
	"path"
)

// stores pointers, assumes they aren't modified anywhere else
type PutUserResponse = db.User

// Handler for the Users resource
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

	form := forms.NewUserEditForm()
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

	fmt.Println("")
	if form.Img != nil {
		err = imaging.Save(form.Img, path.Join(settings.UploadsPath, u.Profile.Img.String))
		if HandleErrForward(w, r, formats.ErrSavingImg, err) != nil {
			return
		}
	}
	if (form.ImgBase64.String == "" || form.Img != nil) && oldImg.Valid {
		_ = os.Remove(path.Join(settings.UploadsPath, oldImg.String))
	}

	if HandleErrForward(w, r, formats.ErrSignupAuthFailure, Authorize(w, r, u)) != nil {
		return
	}

	Handle200(w, r, u)
}
