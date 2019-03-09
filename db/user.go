package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Class representing a user - a JOIN on Account and Profile
// Must not be created or modified directly but im still unsure whether unexported
// fields work with JSON, and how hard my life will be if not.
type User struct {
	Account *Account
	Profile *Profile
}

type UserData struct {
	Pk         int64        `json:"uid"`
	Username   string       `json:"username"`
	Email      string       `json:"email"`
	FirstName  NullString   `json:"firstName"`
	LastName   NullString   `json:"lastName"`
	HighScore  NullInt64    `json:"highScore"`
	Gender     NullString   `json:"gender"`
	Img        NullString   `json:"img"`
	BirthDate  NullDateTime `json:"birthDate"`
	SignupDate time.Time    `json:"signupDate"`
}

func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(UserData{
		u.Account.Pk, u.Account.Username, u.Account.Email,
		u.Profile.FirstName, u.Profile.LastName,
		u.Profile.HighScore, u.Profile.Gender, u.Profile.Img,
		NullDateTime{u.Profile.BirthDate}, u.Profile.SignupDate,
	})
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
	u := &User{&Account{}, &Profile{}}
	err := _db.QueryRow(`
    SELECT a."uid", a."username", a."email", a."password", 
           p."first_name", p."last_name", p."high_score",
           p."gender", p."img", p."birth_date", p."signup_date"
    FROM account a JOIN profile p ON a.uid = p.uid
    WHERE a."uid"=$1`, uid).Scan(&u.Account.Pk, &u.Account.Username, &u.Account.Email,
		&u.Account.Password, &u.Profile.FirstName, &u.Profile.LastName,
		&u.Profile.HighScore, &u.Profile.Gender, &u.Profile.Img,
		&u.Profile.BirthDate, &u.Profile.SignupDate)

	if err != nil {
		return nil, err
	}
	return u, nil
}

func GetUserByUsernameOrEmail(_db Queryable, username string,
	email string) (*User, error) {
	u := &User{&Account{}, &Profile{}}
	err := _db.QueryRow(`
    SELECT a."uid", a."username", a."email", a."password", 
           p."first_name", p."last_name", p."high_score",
           p."gender", p."img", p."birth_date", p."signup_date"
    FROM account a JOIN profile p ON a.uid = p.uid
    WHERE a."email"=$1 OR a."username" = $2`, // both are constrained to be unique
		email, username).Scan(&u.Account.Pk, &u.Account.Username, &u.Account.Email,
		&u.Account.Password, &u.Profile.FirstName, &u.Profile.LastName,
		&u.Profile.HighScore, &u.Profile.Gender, &u.Profile.Img,
		&u.Profile.BirthDate, &u.Profile.SignupDate)

	if err != nil {
		return nil, err
	}
	return u, nil
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
// UserOrderMap
func GetUserMany(_db Queryable, order []SelectOrder,
	limit int, offset int) (*sql.Rows, error) {
	orderStr, err := makeOrderString(UserOrderMap, order)
	if err != nil {
		return nil, err
	}

	limitStr := makeLimitString(limit)
	offsetStr := makeOffsetString(offset)

	return _db.Query(`SELECT a."uid", a."username", a."email", a."password", 
                           p."first_name", p."last_name", p."high_score",
                    p."gender", p."img", p."birth_date", p."signup_date",
                    count(*) OVER() AS n_users
                    FROM account a JOIN profile p ON a.uid = p.uid ` +
		orderStr + limitStr + offsetStr)

}

// Convenience function provided for getting a user out of sql.Rows.Scan()
// I probably should've implemented an sql.Scanner interface...
// the туду was removed because there isn't a universal scanner - sometimes
// we need count(*), sometimes we don't
func UserFromRow(row Scanner) (*User, int, error) {
	u := &User{&Account{}, &Profile{}}
	var nUsers int
	err := row.Scan(&u.Account.Pk, &u.Account.Username, &u.Account.Email,
		&u.Account.Password, &u.Profile.FirstName, &u.Profile.LastName,
		&u.Profile.HighScore, &u.Profile.Gender, &u.Profile.Img,
		&u.Profile.BirthDate, &u.Profile.SignupDate, &nUsers)

	if err == nil {
		return u, nUsers, nil
	} else {
		u.Profile.Pk = u.Account.Pk
		return nil, 0, err
	}
}

// Lists all fields that GetUserMany supports ordering by
// adding an entry to the map is sufficient for everything to work
var UserOrderMap = map[string]string{
	"HighScore":  "p.high_score",
	"FirstName":  "p.first_name",
	"LastName":   "p.last_name",
	"SignupDate": "p.signup_date",
	"Username":   "a.username",
	"Pk":         "a.uid",
}
