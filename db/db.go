// Package implements interface to Postgresql using db.sql and pq PostgreSQL driver
// There's one class per DB entity, containing all the fields
// Every class contains a Save method, which always takes Queryable or sql.DB
// sql.DB is only taken when the body of amethod
// All classes have Pk attributes, which are set to 0 if the object
// is not represented in the database.
package db

// The file provides interfaces and structs necessary for using the db package

import "database/sql"

// Represents the common interface of sql.Tx and sql.DB
type Queryable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// Represents the common interface of Sql.Row and Sql.Rows
type Scanner interface {
	Scan(...interface{}) error
}

// Represents an ordering to the
type SelectOrder struct {
	Field string
	Desc  bool
}

// Returns "order.Field ASC/DESC", provided for convenience
func (order *SelectOrder) String() string {
	if order.Desc {
		return order.Field + " DESC"
	} else {
		return order.Field + " ASC"
	}
}
