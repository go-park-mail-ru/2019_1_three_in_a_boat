package forms

import (
	"github.com/badoux/checkmail"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
)

type SigninForm struct {
	Name     string        `json:"name"`
	Username db.NullString `json:"-"`
	Email    db.NullString `json:"-"`
	Password string        `json:"password"`
	ok       bool          `json:"-"`
}

func (f *SigninForm) Validate() *Report {
	report := NewReport("name", "password")
	report.Ok = true
	fReport := FieldReport{false, []string{formats.ErrFieldTooShort}}

	if len(f.Name) == 0 {
		report.Ok = false
		report.Fields["password"] = fReport
	}

	if len(f.Password) == 0 {
		report.Ok = false
		report.Fields["email"] = fReport
	}

	if !report.Ok {
		return report
	}

	f.parseName()
	f.ok = report.Ok
	return report
}

func (f *SigninForm) parseName() {
	if err := checkmail.ValidateFormat(f.Name); err != nil {
		f.Username.String = f.Name
		f.Username.Valid = true
	} else {
		f.Email.String = f.Name
		f.Email.Valid = true
	}
}

var UnsuccessfulSigninReport = Report{
	false,
	map[string]FieldReport{
		"name":     {false, []string{formats.ErrInvalidCredentials}},
		"password": {false, []string{formats.ErrInvalidCredentials}},
	},
}
