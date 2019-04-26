package db

import (
	"database/sql"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/settings/shared"
)

// Struct representing an author - data, additional to Profile, supplied only
// for the authors of the project.
type Author struct {
	Pk          int64      `json:"uid"`
	DevInfo     NullString `json:"devInfo"`
	Description NullString `json:"description"`
}

// Data required by the /authors API resource, provided for convenience.
// Incorporates Author Class and Username + Img fields of Account and Profile.
type AuthorData struct {
	Pk          int64      `json:"uid"`
	DevInfo     NullString `json:"devInfo"`
	Description NullString `json:"description"`
	Username    string     `json:"username"`
	FirstName   NullString `json:"name"`
	LastName    NullString `json:"lastName"`
	Img         NullString `json:"img"`
}

// Constructor - takes all the fields, does not immediately check if the uid is
// valid. The Pk check is deferred until Save is called.
func NewAuthor(uid int64, devInfo, description NullString) *Author {
	// integrity is handled when saving by the database
	return &Author{
		uid,
		devInfo,
		description,
	}
}

// doesn't really need a GetAuthor so I'm just not gonna bother

// Saves Author object to the database - the user with a.Pk must already exist
// or be created in the same transaction (FK constraint is initially deferred).
// Should be used to mark certain users as authors and add extra info from
// Author struct
func (a *Author) Save(_db Queryable) error {
	_, err := _db.Exec(`
      INSERT INTO author ("uid", "dev_info", "description")
      VALUES ($1, $2, $3)
      ON CONFLICT (uid) DO
      UPDATE
      SET "dev_info" = $2, "description" = $3`,
		// $1   $2         $3
		a.Pk, a.DevInfo, a.Description)
	return err
}

// Creates an AuthorData object from a Scanner returned by GetAllAuthors
func AuthorDataFromRow(row Scanner) (*AuthorData, error) {
	a := &AuthorData{}
	err := row.Scan(&a.Pk, &a.Username, &a.FirstName, &a.LastName,
		&a.Img, &a.DevInfo, &a.Description)
	if err != nil {
		return nil, err
	}
	if a.Img.String == "" || !a.Img.Valid {
		a.Img.String = settings.DefaultImgName
		a.Img.Valid = true
	}
	return a, nil
}

// Returns all authors from the database. Use AuthorDataFromRow on the returned
// object to retrieve the object from a row. And don't forget to Close the Rows!
func GetAllAuthors(_db Queryable) (*sql.Rows, error) {
	return _db.Query(`SELECT a."uid", a."username", p."first_name", p."last_name",
                           p."img", au."dev_info", au."description"
                    FROM author au
                    JOIN Account a ON au.uid = a.uid
                    JOIN profile p ON a.uid = p.uid`)
}
