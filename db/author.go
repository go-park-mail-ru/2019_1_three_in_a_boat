package db

import "database/sql"

// Class representing an entity author
type Author struct {
	Pk          int64      `json:"uid"`
	DevInfo     NullString `json:"devInfo"`
	Description NullString `json:"description"`
}

// Data required by the get_authors API resource, provided for convenience
// Incorporates Author Class and Username + Img fields of Account and Profile
type AuthorData struct {
	Author   *Author    `json:",inline"`
	Username string     `json:"username"`
	Img      NullString `json:"img"`
}

// returns ad.Author.Pk, provided for convenience
func (ad * AuthorData) GetPk() int64 {
  return ad.Author.Pk
}

// Constructor - takes all the fields, does NOT immediately check if uid is valid
// The Pk check is deferred until Save is called
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
// Should be used to mark certain users as authors and add some information
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

// Created an AuthorData object from a Scanner returned by GetAllAuthors
func AuthorDataFromRow(row Scanner) (*AuthorData, error) {
	a := &AuthorData{Author: &Author{}}
	err := row.Scan(&a.Author.Pk, &a.Username,
		&a.Img, &a.Author.DevInfo, &a.Author.Description)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Returns all authors from the database. Use AuthorDataFromRow on the returned
// object to retrieve the object from a row. And don't forget to Close the Rows!
func GetAllAuthors(_db Queryable) (*sql.Rows, error) {
	return _db.Query(`SELECT a."uid", a."username", p."img",
                           au."dev_info", au."description"
                    FROM author au
                    JOIN Account a ON au.uid = a.uid
                    JOIN profile p ON a.uid = p.uid`)
}
