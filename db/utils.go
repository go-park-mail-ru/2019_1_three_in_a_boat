package db

// The file provides utility functions, none of which are exported

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

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
