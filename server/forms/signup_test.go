package forms

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/db"
)

func TestSignupForm_ValidateEmail(t *testing.T) {
	f := SignupForm{}
	cases := []struct {
		str string
		ok  bool
	}{
		{
			str: "foobar",
			ok:  false,
		},
		{
			str: "foobar@foo.bar",
			ok:  true,
		},
		{
			str: "фубар@фу.бар",
			ok:  false,
		},
	}

	for _, c := range cases {
		f.Email = c.str
		if ok := f.ValidateEmail().Ok; ok != c.ok {
			t.Errorf(
				"%s validation: expceted %t, got %t", c.str, c.ok, ok)
		}
	}
}

func TestSignupForm_ValidatePassword(t *testing.T) {
	f := SignupForm{}
	cases := []struct {
		str string
		ok  bool
	}{
		{
			str: "123456",
			ok:  false,
		},
		{
			str: "123456789",
			ok:  true,
		},
		{
			str: strings.Repeat("0123456789", 20),
			ok:  false,
		},
	}

	for _, c := range cases {
		f.Password = c.str
		if ok := f.ValidatePassword().Ok; ok != c.ok {
			t.Errorf(
				"%s validation: expceted %t, got %t", c.str, c.ok, ok)
		}
	}
}

func TestSignupForm_ValidateBirthDate(t *testing.T) {
	f := SignupForm{}
	cases := []struct {
		str string
		ok  bool
	}{
		{
			str: "0-0-0",
			ok:  true,
		},
		{
			str: "14-14-2019",
			ok:  false,
		},
		{
			str: "foobar",
			ok:  false,
		},
		{
			str: "0-0-0-0",
			ok:  false,
		},
	}

	for _, c := range cases {
		f.BirthDateStr = db.NullString{c.str, true}
		if ok := f.ValidateBirthDate().Ok; ok != c.ok {
			t.Errorf(
				"%s validation: expceted %t, got %t", c.str, c.ok, ok)
		}
	}
}

func TestSignupForm_Validate(t *testing.T) {
	cases := []SignupForm{
		{
			Username:     "foobar",
			Password:     "foobarfoobar",
			Email:        "foobar@foo.bar",
			FirstName:    db.NullString{"Foo", true},
			LastName:     db.NullString{"Bar", true},
			BirthDateStr: db.NullString{"0-0-0", true},
			ok:           true,
		},
		{
			Username:     "foobar",
			Password:     "foobarfoobar",
			Email:        "foobar@foo..bar",
			FirstName:    db.NullString{"Foo", true},
			LastName:     db.NullString{"Bar", true},
			BirthDateStr: db.NullString{"0-0-0", true},
			ok:           false,
		},
		{
			Username:     "foobar",
			Password:     "fo",
			Email:        "foobar@foo.bar",
			FirstName:    db.NullString{"Foo", true},
			LastName:     db.NullString{"Bar", true},
			BirthDateStr: db.NullString{"0-0-0", true},
			ok:           false,
		},
	}

	for _, c := range cases {
		valid := c.ok
		if ok := c.Validate().Ok; ok != valid {
			t.Errorf(
				"expceted %t, got %t", valid, ok)
		}
	}
}

func TestSignupForm_MakeUser(t *testing.T) {
	cases := []struct {
		SignupForm SignupForm
		User       db.User
	}{
		{
			SignupForm: SignupForm{
				Username:     "foobar",
				Password:     "foobarfoobar",
				Email:        "foobar@foo.bar",
				FirstName:    db.NullString{"Foo", true},
				BirthDateStr: db.NullString{"0-0-0", true},
				ok:           true,
			},
			User: db.User{
				Account: &db.Account{
					Pk:       0,
					Username: "foobar",
					Email:    "foobar@foo.bar",
				},
				Profile: &db.Profile{
					Pk:        0,
					FirstName: db.NullString{"Foo", true},
					LastName:  db.NullString{"", false},
					BirthDate: db.NullTime{Valid: false},
					HighScore: db.NullInt64{0, false},
					Gender:    db.NullString{"", false},
					Img:       db.NullString{"", false},
				},
			},
		},
		{
			SignupForm: SignupForm{
				ok: false,
			},
		},
	}

	for _, c := range cases {
		u, err := c.SignupForm.MakeUser()
		if c.SignupForm.ok {
			pwdMatch, pwdErr :=
				db.AccountComparePasswordToHash(c.SignupForm.Password, u.Account.Password)
			if pwdErr != nil {
				t.Errorf("hashing error")
			} else if err != nil {
				t.Errorf("epxceted a user, got an error")
			} else if u.Account.Pk != c.User.Account.Pk ||
				u.Account.Username != c.User.Account.Username ||
				u.Account.Email != c.User.Account.Email ||
				!pwdMatch ||
				u.Profile.FirstName != c.User.Profile.FirstName ||
				u.Profile.LastName != c.User.Profile.LastName ||
				u.Profile.BirthDate != c.User.Profile.BirthDate ||
				u.Profile.HighScore != c.User.Profile.HighScore ||
				u.Profile.Gender != c.User.Profile.Gender ||
				u.Profile.Img != c.User.Profile.Img {
				fmt.Println(u.Account.Pk != c.User.Account.Pk,
					u.Account.Username != c.User.Account.Username,
					u.Account.Email != c.User.Account.Email,
					pwdMatch,
					u.Profile.FirstName != c.User.Profile.FirstName,
					u.Profile.LastName != c.User.Profile.LastName,
					u.Profile.BirthDate != c.User.Profile.BirthDate,
					u.Profile.HighScore != c.User.Profile.HighScore,
					u.Profile.Gender != c.User.Profile.Gender,
					u.Profile.Img != c.User.Profile.Img)
				t.Errorf("corrupted/empty user data received")
			}
		}

	}
}
