// Package implements interface to Postgresql using db.sql and pq PostgreSQL driver
// There's one class per DB entity, containing all the fields
// Every class contains a Save method, which always takes Queryable or sql.DB
// sql.DB is only taken when the body of amethod
// All classes have Pk attributes, which are set to 0 if the object
// is not represented in the database.
package db

// The file provides utility functions, structs and interfaces.
// All of the structs and interfaces are exported. None of the functions are.

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// represents the common interface of sql.Tx and sql.DB
type Queryable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// represents the common interface of Sql.Row and Sql.Rows
type Scanner interface {
	Scan(...interface{}) error
}

// represents an ordering to the
type SelectOrder struct {
	Field string
	Desc  bool
}

// returns "order.Field ASC/DESC", provided for convenience
func (order *SelectOrder) String() string {
	if order.Desc {
		return order.Field + " DESC"
	} else {
		return order.Field + " ASC"
	}
}

// Makes a valid SQL ORDER BY statement out of a slice of SelectOrder
// adds a space at the end, so chaining statements with + is safe
func makeOrderString(orderMap map[string]string, order []SelectOrder) (string, error) {
	if len(order) == 0 {
		return "", nil
	}

	orderBuilder := strings.Builder{}
	orderBuilder.WriteString("ORDER BY ")
	for i, orderElt := range order {
		if val, ok := orderMap[orderElt.Field]; ok {
			orderBuilder.WriteString(val)
			if orderElt.Desc {
				orderBuilder.WriteString(" DESC")
			}
			if i != len(order)-1 {
				orderBuilder.WriteString(", ")
			} else {
				orderBuilder.WriteString(" ")
			}
		} else {
			return "", errors.New(fmt.Sprintf(
				"cannot order by %s: invalid field", orderElt.Field))
		}
	}

	return orderBuilder.String(), nil

}

// Makes a valid SQL LIMIT statement based on provided limit.
// If limit is negative, returns an empty string
func makeLimitString(limit int) string {
	if limit < 0 {
		return ""
	} else {
		return fmt.Sprintf("LIMIT %d ", limit)
	}
}

// Makes a valid SQL OFFSET statement based on provided limit.
// If limit is negative or zero, returns an empty string
func makeOffsetString(offset int) string {
	if offset <= 0 {
		return ""
	} else {
		return fmt.Sprintf("OFFSET %d ", offset)
	}
}

// If err is not nil aborts transaction. Returns original error, transaction error
func abortOnError(tx *sql.Tx, _err error) (err, txError error) {
	if _err != nil {
		return _err, tx.Rollback()
	}
	return nil, nil
}

// If err is not nil aborts transaction, otherwise commits it. Return values are the same
// as in abortOnError
func abortOnErrorOrCommit(tx *sql.Tx, _err error) (err, txError error) {
	if _err != nil {
		return _err, tx.Rollback()
	} else {
		return nil, tx.Commit()
	}
}
