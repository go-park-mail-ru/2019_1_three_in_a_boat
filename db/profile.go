package db

import (
	"time"
)

// Class representing a profile - all the non-essential data of a user
// Must not be created or modified directly but im still unsure whether unexported
// fields work with JSON, and how hard my life will be if not.
type Profile struct {
	Pk         int64      `json:"-"`
	FirstName  NullString `json:"firstName"`
	LastName   NullString `json:"lastName"`
	HighScore  NullInt64  `json:"highScore"`
	Gender     NullString `json:"gender"`
	Img        NullString `json:"img"`
	BirthDate  NullTime   `json:"birthDate"`
	SignupDate time.Time  `json:"signupDate"`
}

// Constructor - takes all the fields, does NOT immediately check if uid is valid
// The Pk check is deferred until Save is called
func NewProfile(uid int64, firstName, lastName NullString,
	highScore NullInt64, gender, img NullString,
	birthDate NullTime) (*Profile, error) {
	return &Profile{uid,
		firstName,
		lastName,
		highScore,
		gender,
		img,
		birthDate,
		time.Time{},
	}, nil
}

// Retrieves a single profile from the database, forwards the SQL error, if any
func GetProfile(_db Queryable, uid int64) (*Profile, error) {
	p := &Profile{}
	err := _db.QueryRow(`SELECT first_name, last_name, high_score, gender, 
                              img, birth_date, signup_date
                       FROM profile WHERE uid = $1`, uid).Scan(
		p.FirstName, p.LastName, p.HighScore, p.Gender, p.Img, p.BirthDate, p.SignupDate)
	if err != nil {
		return nil, err
	} else {
		return p, err
	}
}

// Saves User object to the database - the user with p.Pk must already exist
// or be created in the same transaction (FK constraint is initially deferred).
// Should not be used to create new users - use User class instead.
func (p *Profile) Save(_db Queryable) error {
	_, err := _db.Exec(`
       INSERT INTO profile (uid, first_name, last_name, high_score,
                            gender, img, birth_date) 
       VALUES ($1, $2, $3, $4, $5, $6, $7)
		   ON CONFLICT (uid) DO 
		   UPDATE
       SET first_name = $2, last_name = $3, high_score = $4,
           gender = $5,  img = $6, birth_date = $7`,
		// $1   $2           $3          $4           $5        $6     $7
		p.Pk, p.FirstName, p.LastName, p.HighScore, p.Gender, p.Img, p.BirthDate)
	return err
}
