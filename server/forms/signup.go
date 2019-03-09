package forms

import (
	"errors"
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"github.com/lib/pq"
	"github.com/nbutton23/zxcvbn-go"
	"strings"
	"time"
)

type SignupForm struct {
	Username     string        `json:"username"`
	Password     string        `json:"password"`
	Email        string        `json:"email"`
	FirstName    db.NullString `json:"name,omitempty"`
	LastName     db.NullString `json:"lastname,omitempty"`
	BirthDateStr string        `json:"date,omitempty"`
	BirthDate    db.NullTime   `json:"-"`
	ok           bool          `json:"-"`
}

func (f *SignupForm) Validate() *Report {
	report := NewReport()

	report.Fields["username"] = f.ValidateUsername()
	report.Fields["password"] = f.ValidatePassword()
	report.Fields["email"] = f.ValidateEmail()
	report.Fields["name"] = f.ValidateFirstName()
	report.Fields["lastname"] = f.ValidateLastName()
	report.Fields["date"] = f.ValidateLastName()

	// optional fields return OK if they're empty - error only if the data is
	// there, but it's invalid
	report.Ok = report.Fields["username"].Ok && report.Fields["password"].Ok &&
		report.Fields["email"].Ok && report.Fields["name"].Ok &&
		report.Fields["lastname"].Ok && report.Fields["date"].Ok
	f.ok = report.Ok

	return report
}

func (f *SignupForm) MakeUser() (*db.User, error) {
	if !f.ok {
		return nil, errors.New("can not create a user from an invalid form")
	}
	a, err := db.NewAccount(f.Username, f.Email, f.Password)
	if err != nil {
		return nil, err
	}

	p, err := db.NewProfile(0, f.FirstName, f.LastName,
		db.NullInt64{0, false}, db.NullString{"", false},
		db.NullString{"", false}, f.BirthDate)
	if err != nil {
		return nil, err
	}

	u, err := db.NewUser(a, p)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (f *SignupForm) ValidateUsername() (r FieldReport) {
	f.Username = strings.TrimSpace(f.Username)
	return checkLength(f.Username, 3, 32)
}

func (f *SignupForm) ValidatePassword() (r FieldReport) {
	if len([]byte(f.Password)) > 128 { // What if they're trying to overload our hash?
		r.Errors = []string{formats.ErrFieldTooLong}
	} else if len([]byte(f.Password)) < 8 {
		r.Errors = []string{formats.ErrFieldTooShort}
	}

	strength := zxcvbn.PasswordStrength(f.Password, []string{f.Username, f.Email,
		f.FirstName.String, f.LastName.String})

	if strength.Score < minPasswordStrnegth {
		r.Errors = append(r.Errors, formats.ErrPasswordTooWeak[strength.Score])
	}

	r.Ok = len(r.Errors) == 0
	return
}

func (f *SignupForm) ValidateEmail() (r FieldReport) {
	f.Email = strings.TrimSpace(f.Email)
	if err := checkmail.ValidateFormat(f.Email); err != nil {
		r.Errors = []string{formats.ErrInvalidEmail}
		return
	}

	if !emailExistsCheck {
		r.Ok = true
		return
	}

	if err := checkmail.ValidateHost(f.Email); err != nil {
		r.Errors = []string{formats.ErrEmailDomainDoesNotExist}
	} else {
		err := checkmail.ValidateHost(f.Email)
		if _, valid := err.(checkmail.SmtpError); valid && err != nil {
			r.Errors = []string{formats.ErrEmailDoesNotExist}
		} else {
			r.Ok = true
		}
	}

	return
}

func (f *SignupForm) ValidateBirthDate() (r FieldReport) {
	// BirthDateStr is generated by JS so no point in trimming
	var err error
	if f.BirthDateStr == "0-0-0" {
		// f.BirthDate stays Null
		return FieldReport{true, []string{}}
	}

	f.BirthDate.Time, err = time.Parse(dateFormat, f.BirthDateStr)
	if err != nil {
		r.Errors = append(r.Errors, formats.ErrInvalidDate)
	} else if f.BirthDate.Time.Year() < 1900 || time.Now().After(f.BirthDate.Time) {
		r.Errors = append(r.Errors, formats.ErrDateOutOfRange)
	} else {
		r.Ok = true
		f.BirthDate.Valid = true
	}

	return
}

func (f *SignupForm) ValidateFirstName() (r FieldReport) {
	f.FirstName.String = strings.TrimSpace(f.FirstName.String)
	return checkLengthOptional(&f.LastName, 1, 32)
}

func (f *SignupForm) ValidateLastName() (r FieldReport) {
	f.LastName.String = strings.TrimSpace(f.LastName.String)
	return checkLengthOptional(&f.LastName, 1, 32)
}

func CheckUserDbConstraints(err error) (*Report, error) {
	if err != nil {
		report := NewReport("username", "password", "email",
			"name", "lastname", "date")
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Constraint != "" {
			fieldReport :=
				FieldReport{false, []string{formats.ErrUniqueViolation}}
			fmt.Println(pqErr)

			if pqErr.Constraint == "account_username_key" {
				report.Fields["username"] = fieldReport
			} else if pqErr.Constraint == "account_email_key" {
				report.Fields["password"] = fieldReport
			} // else shouldn't ever be possible
			return report, errors.New(formats.ErrUniqueViolation)
		}
		return nil, err
	}

	return nil, nil
}
