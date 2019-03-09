package forms

import (
	"bytes"
	"encoding/base64"
	"github.com/disintegration/imaging"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/formats"
	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/server/settings"
	"github.com/google/uuid"
	"image"
	"strings"
)

type UserEditForm struct {
	*SignupForm
	ImgBase64 db.NullString `json:"img"` // the one that's read from JSON
	Img       image.Image   `json:"-"`   // resized img to be conveniently saved
	ImgName   string        `json:"-"`   // path to the file, not accounting for settings
	Gender    db.NullString `json:"gender"`
	ok        bool          `json:"-"`
}

func NewUserEditForm() *UserEditForm {
	return &UserEditForm{SignupForm: &SignupForm{}}
}

func (f *UserEditForm) Validate() *Report {
	report := NewReport()

	report.Fields["username"] = f.ValidateUsername()
	report.Fields["password"] = f.ValidatePassword()
	report.Fields["email"] = f.ValidateEmail()
	report.Fields["name"] = f.ValidateFirstName()
	report.Fields["lastname"] = f.ValidateLastName()
	report.Fields["date"] = f.ValidateLastName()
	report.Fields["img"] = f.ValidateImg()

	report.Ok = report.Fields["username"].Ok && report.Fields["password"].Ok &&
		report.Fields["email"].Ok && report.Fields["name"].Ok &&
		report.Fields["lastname"].Ok && report.Fields["date"].Ok &&
		report.Fields["img"].Ok
	f.ok = report.Ok
	return report
}

func (f *UserEditForm) ValidateImg() FieldReport {
	if !f.ImgBase64.Valid {
		f.ImgBase64.Valid = false
		return FieldReport{true, nil}
	}

	if f.ImgBase64.String == "" { // !Valid already checked
		f.ImgBase64.Valid = true
		return FieldReport{true, nil}
	}

	imgBytes, err := base64.StdEncoding.DecodeString(f.ImgBase64.String)
	if err != nil {
		return FieldReport{false, []string{formats.ErrBase64Decoding}}
	}

	img, err := imaging.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return FieldReport{false, []string{formats.ErrBase64Decoding}}
	}

	f.Img = imaging.Fill(img, settings.ImageSize[0], settings.ImageSize[1],
		imaging.Center, imaging.Lanczos)
	f.ImgName = uuid.New().String() + ".jpg"

	return FieldReport{true, nil}
}

func (f *UserEditForm) ValidatePassword() FieldReport {
	f.Password = strings.TrimSpace(f.Password)
	if f.Password != "" {
		return f.SignupForm.ValidatePassword()
	} else {
		return FieldReport{true, nil}
	}
}

func (f *UserEditForm) ValidateUsername() FieldReport {
	f.Username = strings.TrimSpace(f.Username)
	if f.Username != "" {
		return f.SignupForm.ValidateUsername()
	} else {
		return FieldReport{true, []string{}}
	}
}

func (f *UserEditForm) ValidateEmail() FieldReport {
	f.Email = strings.TrimSpace(f.Email)
	if f.Email != "" {
		return f.SignupForm.ValidateEmail()
	} else {
		return FieldReport{true, []string{}}
	}
}

func (f *UserEditForm) ValidateGender() FieldReport {
	f.Gender.String = strings.TrimSpace(f.Gender.String)
	if f.Gender.String != "" && f.Gender.Valid {
		switch f.Gender.String {
		case "male":
			fallthrough
		case "female":
			fallthrough
		case "other":
			return FieldReport{true, []string{}}
		default:
			return FieldReport{false, []string{}}
		}
	} else {
		f.Gender.Valid = false
		return FieldReport{true, []string{}}
	}
}

func (f *UserEditForm) EditUser(u *db.User) (*db.User, error) {
	if !f.ok {
		return nil, ErrFormInvalid
	}

	var err error
	if f.Password != "" {
		u.Account.Password, err = db.AccountGeneratePasswordHash(f.Password)
		if err != nil {
			return nil, err
		}
	}

	if f.Username != "" {
		u.Account.Username = f.Username
	}

	if f.Email != "" {
		u.Account.Email = f.Email
	}

	if f.ImgBase64.Valid {
		u.Profile.Img.Valid = true
		u.Profile.Img.String = f.ImgName
	} else if f.ImgBase64.String == "" {
		u.Profile.Img.Valid = false
		u.Profile.Img.String = ""
	} // else (null) do nothing

	u.Profile.FirstName = f.FirstName
	u.Profile.LastName = f.LastName
	u.Profile.Gender = f.Gender
	u.Profile.BirthDate = f.BirthDate

	return u, nil
}
