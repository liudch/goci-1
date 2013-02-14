package goci

import (
	"database/sql/driver"
)

type statement struct{}

func (stmt *statement) Close() (err error) {
	return nil
}

func (stmt *statement) NumInput() (num int) {
	return 0
}

// Exec executes a query that doesn't return rows, such
// as an INSERT or UPDATE.
func (stmt *statement) Exec(v []driver.Value) (result driver.Result, err error) {
	return nil, nil
}

// Exec executes a query that may return rows, such as a
// SELECT.
func (stmt *statement) Query(v []driver.Value) (_ driver.Rows, err error) {
	return nil, nil
}
