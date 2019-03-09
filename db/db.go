// Package implements interface to Postgresql using db.sql and pq PostgreSQL
// driver. There's one class per DB entity, containing all the fields from DB.
// Every class has a Save method, which always takes Queryable or sql.DB. sql.DB
// is only required when the body of a method requires execution in a separate
// transaction. All classes have Pk attributes, which are set to 0 if the object
// is not represented in the database, and to its primary key if it is.
package db

// The file provides interfaces and structs necessary for using the db package

import "database/sql"

// Represents the common interface of sql.Tx and sql.DB, accepted by most of the
// classes that require a connection.
type Queryable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// Represents the common interface of Sql.Row and Sql.Rows, returned by the
// Get* methods.
type Scanner interface {
	Scan(...interface{}) error
}

// Represents an ordering accepted by Get*Many methods that support ordering
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
