package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// Class representing a user - a JOIN on Account and Profile
// Must not be created or modified directly but im still unsure whether unexported
// fields work with JSON, and how hard my life will be if not.
type User struct {
	Account *Account `json:",inline"`
	Profile *Profile `json:",inline"`
}

// returns u.Account.Pk, provided for convenience
func (u *User) GetPk() int64 {
	return u.Account.Pk
}

// Constructor - takes Account, constructed by NewAccount and profile,
// constructed by NewProfile. If (and only if) PKs do not coincide, returns an error.
func NewUser(account *Account, profile *Profile) (*User, error) {
	if account.Pk != profile.Pk {
		return nil, errors.New(fmt.Sprintf(
			"account.Pk != profile.Pk: %d != %d", account.Pk, profile.Pk))
	}

	return &User{account, profile}, nil
}

// Retrieves a single user from the database, forwards the SQL error, if any
func GetUser(_db Queryable, uid int64) (*User, error) {
	// TODO: write a join, jeez what is this shittery
	a, err := GetAccount(_db, uid)
	if err != nil {
		return nil, err
	}

	p, err := GetProfile(_db, uid)
	if err != nil {
		return nil, err
	}

	return &User{a, p}, nil
}

// Saves User object to the database -
// no matter if it's a newly created or an existing object
// only accepts sql.DB, since it needs to use its own transaction
func (u *User) Save(_db *sql.DB) (err, transactionError error) {
	tx, transactionError := _db.Begin()
	if transactionError != nil {
		return
	}

	err = u.Account.Save(tx)
	if e, txE := abortOnError(tx, err); e != nil || txE != nil {
		return e, txE
	}
	u.Profile.Pk = u.Account.Pk

	err = u.Profile.Save(tx)
	return abortOnErrorOrCommit(tx, err)
}

// Retrieves ALL values from multiple rows of the account/profile join
// does not support projection, supports ordering on all fields specified in
// userOrderMap
func GetUserMany(_db Queryable, order []SelectOrder,
	limit int, offset int) (*sql.Rows, error) {
	orderStr, err := makeOrderString(userOrderMap, order)
	if err != nil {
		return nil, err
	}

	limitStr := makeLimitString(limit)
	offsetStr := makeOffsetString(offset)

	return _db.Query(`SELECT a."uid", a."username", a."email", a."password", 
                           p."first_name", p."last_name", p."high_score", p."gender", 
                           p."img", p."birth_date", p."signup_date"
                    FROM account a JOIN profile p ON a.uid = p.uid ` +
		orderStr + limitStr + offsetStr)

}

// Convenience function provided for getting a user out of sql.Rows.Scan()
// I probably should've implemented an sql.Scanner interface...
// TODO: implement a Scanner interface for User
func UserFromRow(row Scanner) (*User, error) {
	u := &User{&Account{}, &Profile{}}
	err := row.Scan(&u.Account.Pk, &u.Account.Username, &u.Account.Email,
		&u.Account.Password, &u.Profile.FirstName, &u.Profile.LastName, &u.Profile.HighScore,
		&u.Profile.Gender, &u.Profile.Img, &u.Profile.BirthDate, &u.Profile.SignupDate)

	if err == nil {
		return u, nil
	} else {
		u.Profile.Pk = u.Account.Pk
		return nil, err
	}
}

// Lists all fields that GetUserMany supports ordering by
// adding an entry to the map is sufficient for everything to work
var userOrderMap = map[string]string{
	"HighScore":  "p.high_score",
	"FirstName":  "p.first_name",
	"LastName":   "p.last_name",
	"SignupDate": "p.signup_date",
	"Username":   "a.username",
	"Pk":         "a.uid",
}
