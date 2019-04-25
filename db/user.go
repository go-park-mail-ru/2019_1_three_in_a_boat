package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/settings/shared"
)

// Class representing a user - a JOIN on Account and Profile
// Must not be created or modified directly but im still unsure whether unexported
// fields work with JSON, and how hard my life will be if not.
type User struct {
	Account *Account
	Profile *Profile
}

// Helper class, representing the User in a way it is serailized. User JSON
// interface converts User to UserData and forwards the marshal/unmarshal call.
type UserData struct {
	Pk        int64      `json:"uid"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	HighScore NullInt64  `json:"highScore"`
	Gender    NullString `json:"gender"`
	Img       NullString `json:"img"`
}

func (u User) MarshalJSON() ([]byte, error) {
	img := u.Profile.Img
	if img.String == "" || img.Valid == false {
		img.String = settings.DefaultImgName
		img.Valid = true
	}
	return json.Marshal(UserData{
		Pk:        u.Account.Pk,
		Username:  u.Account.Username,
		Email:     u.Account.Email,
		HighScore: u.Profile.HighScore,
		Gender:    u.Profile.Gender,
		Img:       img,
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

// Retrieves a single user from the database, if either username or email
// are found in the DB. To make the response deterministic, one string should
// be empty, unless you're absolutely sure both belong to the same User.
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

// Saves a User in the database, forwards the SQL error, if any. Can be used for
// both updating and inserting. The operation is determined based on Pk: 0 ->
// INSERT, non-0 -> update. Return error can be examined via type assertion to
// the underlying driver's error type
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
// returned by GetUserMany. For single users, use GetUser or
// GetUserByNameOrEmail User does not implement Scanner interface because there
// isn't a universal scanner - sometimes we need count(*), sometimes we don't.
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

func UpdateScoreById(_db Queryable, uid int64, score int64) error {
	_, err := _db.Exec(
		`UPDATE profile p SET high_score=$2 WHERE p."uid"=$1 AND
           (p."high_score" < $2 OR p."high_score" IS NULL)`,
		uid, score)
	return err
}

// Lists all fields that GetUserMany supports ordering by. Adding an entry to
// this map is sufficient for ordering with the field to work
var UserOrderMap = map[string]string{
	"Highscore":   "p.high_score",
	"First_name":  "p.first_name",
	"Last_name":   "p.last_name",
	"Signup_date": "p.signup_date",
	"Username":    "a.username",
	"Pk":          "a.uid",
}
