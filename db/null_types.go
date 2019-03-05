package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"time"
)

type NullString sql.NullString

func (ns NullString) MarshallJSON() ([]byte, error) {
	if val, err := ns.Value(); err != nil {
		return nil, err
	} else {
		return json.Marshal(val)
	}
}

func (ns NullString) Value() (driver.Value, error) {
	return sql.NullString(ns).Value()
}

func (ns *NullString) Scan(value interface{}) error {
	return (*sql.NullString)(ns).Scan(value)
}

type NullInt64 sql.NullInt64

func (ni NullInt64) MarshallJSON() ([]byte, error) {
	if val, err := ni.Value(); err != nil {
		return nil, err
	} else {
		return json.Marshal(val)
	}
}

func (ni *NullInt64) Scan(value interface{}) error {
	return (*sql.NullInt64)(ni).Scan(value)
}

func (ni NullInt64) Value() (driver.Value, error) {
	return sql.NullInt64(ni).Value()
}

type NullFloat64 sql.NullFloat64

func (nf NullFloat64) MarshallJSON() ([]byte, error) {
	if val, err := nf.Value(); err != nil {
		return nil, err
	} else {
		return json.Marshal(val)
	}
}

func (nf *NullFloat64) Scan(value interface{}) error {
	return (*sql.NullFloat64)(nf).Scan(value)
}

func (nf NullFloat64) Value() (driver.Value, error) {
	return sql.NullFloat64(nf).Value()
}

type NullBool sql.NullBool

func (nb NullBool) MarshallJSON() ([]byte, error) {
	if val, err := nb.Value(); err != nil {
		return nil, err
	} else {
		return json.Marshal(val)
	}
}

func (nb *NullBool) Scan(value interface{}) error {
	return (*sql.NullBool)(nb).Scan(value)
}

func (nb NullBool) Value() (driver.Value, error) {
	return sql.NullBool(nb).Value()
}

type NullTime pq.NullTime

func (nt NullTime) MarshallJSON() ([]byte, error) {
	if val, err := nt.Value(); err != nil {
		return nil, err
	} else {
		return json.Marshal(val)
	}
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
