package db

import "fmt"

type ObjectNotFoundError struct {
	requestedId uint64
}

func (e *ObjectNotFoundError) Error() string {
	return fmt.Sprintf("cannot find %d in the database", e.requestedId)
}

type TableNotFoundError struct {
	table string
}

func (e *TableNotFoundError) Error() string {
	return fmt.Sprintf("cannot find %s in the database", e.table)
}