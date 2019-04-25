package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
)

const DateFormat = "02.01.2006"

// Represents a NULLable string from the DB. Implements JSON and DB IO
// interfaces.
type NullString sql.NullString

func (ns NullString) MarshalJSON() ([]byte, error) {
	if val, err := ns.Value(); err != nil {
		return nil, err
	} else {
		return json.Marshal(val)
	}
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var i *string
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	if i != nil {
		ns.Valid = true
		ns.String = *i
	} else {
		ns.Valid = false
	}

	return nil
}

func (ns NullString) Value() (driver.Value, error) {
	return sql.NullString(ns).Value()
}

func (ns *NullString) Scan(value interface{}) error {
	return (*sql.NullString)(ns).Scan(value)
}

// Represents a NULLable Int64 from the DB. Implements JSON and DB IO
// interfaces.
type NullInt64 sql.NullInt64

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if val, err := ni.Value(); err != nil {
		return nil, err
	} else {
		return json.Marshal(val)
	}
}

func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	var i *int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	if i != nil {
		ni.Valid = true
		ni.Int64 = *i
	} else {
		ni.Valid = false
	}

	return nil
}

func (ni *NullInt64) Scan(value interface{}) error {
	return (*sql.NullInt64)(ni).Scan(value)
}

func (ni NullInt64) Value() (driver.Value, error) {
	return sql.NullInt64(ni).Value()
}

// Represents a NULLable Float64 from the DB. Implements JSON and DB IO
// interfaces.
type NullFloat64 sql.NullFloat64

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if val, err := nf.Value(); err != nil {
		return nil, err
	} else {
		return json.Marshal(val)
	}
}

func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
	var f *float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}

	if f != nil {
		nf.Valid = true
		nf.Float64 = *f
	} else {
		nf.Valid = false
	}

	return nil
}

func (nf *NullFloat64) Scan(value interface{}) error {
	return (*sql.NullFloat64)(nf).Scan(value)
}

func (nf NullFloat64) Value() (driver.Value, error) {
	return sql.NullFloat64(nf).Value()
}

// Represents a NULLable Bool from the DB. Implements JSON and DB IO
// interfaces.
type NullBool sql.NullBool

func (nb NullBool) MarshalJSON() ([]byte, error) {
	if val, err := nb.Value(); err != nil {
		return nil, err
	} else {
		return json.Marshal(val)
	}
}

func (nb *NullBool) UnmarshalJSON(data []byte) error {
	var b *bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}

	if b != nil {
		nb.Valid = true
		nb.Bool = *b
	} else {
		nb.Valid = false
	}

	return nil
}

func (nb *NullBool) Scan(value interface{}) error {
	return (*sql.NullBool)(nb).Scan(value)
}

func (nb NullBool) Value() (driver.Value, error) {
	return sql.NullBool(nb).Value()
}

// Represents a NULLable Time from the DB. Implements JSON and DB IO
// interfaces.
type NullTime pq.NullTime

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if val, err := nt.Value(); err != nil {
		return nil, err
	} else {
		return json.Marshal(val)
	}
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var i *time.Time
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	if i != nil {
		nt.Valid = true
		nt.Time = *i
	} else {
		nt.Valid = false
	}

	return nil
}

func (nt *NullTime) Scan(value interface{}) error {
	// the pq implementation is bad
	switch value.(type) {
	case time.Time:
		nt.Time, nt.Valid = value.(time.Time), true
	case nil:
		nt.Time, nt.Valid = time.Time{}, false
	default:
		nt.Time, nt.Valid = time.Time{}, false
		return fmt.Errorf("sql: db: converting driver.Value type "+
			"%T to a db.NullTime", value)
	}
	return nil
}

func (nt NullTime) Value() (driver.Value, error) {
	return pq.NullTime(nt).Value()
}

// Represents a NULLable Date from the DB. Implements JSON and DB IO
// interfaces. The difference between this and NullTime is that NullDateTime
// serializes into a DD.MM.YYYY format, as defined in the DateFormat variable
type NullDateTime struct {
	NullTime
}

func (ndt NullDateTime) MarshalJSON() ([]byte, error) {
	if !ndt.Valid {
		return ndt.NullTime.MarshalJSON()
	} else {
		return NullString{
			String: ndt.NullTime.Time.Format(DateFormat),
			Valid:  true,
		}.MarshalJSON()
	}
}

func (ndt *NullDateTime) UnmarshalJSON(data []byte) error {
	var nullStr NullString
	var err error
	if err = json.Unmarshal(data, &nullStr); err != nil {
		return err
	}
	if nullStr.Valid {
		if ndt.Time, err = time.Parse(DateFormat, nullStr.String); err != nil {
			return err
		} else {
			ndt.Valid = true
		}
	} else {
		ndt.Time = time.Time{}
		ndt.Valid = false
	}

	return nil
}

func (ndt *NullDateTime) Scan(value interface{}) error {
	return ndt.NullTime.Scan(value)
}

func (ndt NullDateTime) Value() (driver.Value, error) {
	return pq.NullTime(ndt.NullTime).Value()
}
