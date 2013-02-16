package goci

/*
#include <oci.h>
#include <stdlib.h>
#include <string.h>

#cgo pkg-config: oci8
*/
import "C"
import (
	"database/sql/driver"
	"io"
  "fmt"
)

type rows struct {
	stmt    *statement
	columns []column
}

// Columns returns the names of the columns. The number of
// columns of the result is inferred from the length of the
// slice.  If a particular column name isn't known, an empty
// string should be returned for that entry.
func (r *rows) Columns() []string {
	names := make([]string, len(r.columns))
	for pos, column := range r.columns {
		names[pos] = column.name
	}
	return names

}

// Close closes the rows iterator.
func (r *rows) Close() error {
  return r.stmt.Close()
}

// Next is called to populate the next row of data into
// the provided slice. The provided slice will be the same
// size as the Columns() are wide.
//
// The dest slice may be populated only with
// a driver Value type, but excluding string.
// All string values must be converted to []byte.
//
// Next should return io.EOF when there are no more rows.
func (r *rows) Next(dest []driver.Value) error {
	rv := C.OCIStmtFetch2((*C.OCIStmt)(r.stmt.handle), (*C.OCIError)(r.stmt.conn.err), 1, C.OCI_FETCH_NEXT, 1, C.OCI_DEFAULT)
  if rv == C.OCI_ERROR {
    err :=  ociGetError(r.stmt.conn.err)
    fmt.Println(err)
    return ociGetError(r.stmt.conn.err)
  }
  if rv == C.OCI_NO_DATA {
    fmt.Println("no data")
    return io.EOF
  }
	
  for i := range dest {
    dest[i] = string(r.columns[i].raw)
  }
	return nil
}
