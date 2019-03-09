package forms

import (
	"github.com/badoux/checkmail"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
)

// A form responsible for validating signin form. Only checks validity of the
// fields, does work with the database, that is done in the corresponding route.
// Implicitly implements JSON Unmarshaller interface.
type SigninForm struct {
	Name     string        `json:"name"`
	Username db.NullString `json:"-"`
	Email    db.NullString `json:"-"`
	Password string        `json:"password"`
	ok       bool          `json:"-"`
}

// Returns a report indicating whether the form is valid. If report.Ok, the
// only one of the form.Username and form.Email will be non-empty, so it's safe
// to send it to db.GetUserByUsernameOrEmail. Checks only that the fields are
// non-empty.
func (f *SigninForm) Validate() *Report {
	report := NewReport("name", "password")
	report.Ok = true
	fReport := FieldReport{false, []string{formats.ErrFieldTooShort}}

	if len(f.Name) == 0 {
		report.Ok = false
		report.Fields["name"] = fReport
	}

	if len(f.Password) == 0 {
		report.Ok = false
		report.Fields["password"] = fReport
	}

	if !report.Ok {
		return report
	}

	f.parseName()
	f.ok = report.Ok
	return report
}

// sets f.Username or f.Email to f.Name, depending on whether f.Name is a valid email
func (f *SigninForm) parseName() {
	if err := checkmail.ValidateFormat(f.Name); err != nil {
		f.Username.String = f.Name
		f.Username.Valid = true
	} else {
		f.Email.String = f.Name
		f.Email.Valid = true
	}
}

// A report that is never used, in the package and provided to the route as a
// report that should be sent to the client. It does not indicate whether the
// user wasn't found or the passwords didn't match.
var UnsuccessfulSigninReport = Report{
	false,
	map[string]FieldReport{
		"name":     {false, []string{formats.ErrInvalidCredentials}},
		"password": {false, []string{formats.ErrInvalidCredentials}},
	},
}
