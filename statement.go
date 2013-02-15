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
	"fmt"
	"unsafe"
)

type statement struct {
	handle unsafe.Pointer
	conn   *connection
}

func (stmt *statement) Close() error {
	return nil
}

func (stmt *statement) NumInput() int {
	return 0
}

// Exec executes a query that doesn't return rows, such
// as an INSERT or UPDATE.
func (stmt *statement) Exec(v []driver.Value) (driver.Result, error) {
	return nil, nil
}

// Exec executes a query that may return rows, such as a
// SELECT.
func (stmt *statement) Query(v []driver.Value) (driver.Rows, error) {
	if err := stmt.bind(v); err != nil {
		return nil, err
	}

	return nil, nil
}

func (stmt *statement) bind(args []driver.Value) error {
	var binding *C.OCIBind
	for pos, value := range args {
		buffer := []byte(fmt.Sprint("%v", value))
		buffer = append(buffer, 0)
		result := C.OCIBindByPos((*C.OCIStmt)(stmt.handle), &binding, (*C.OCIError)(stmt.conn.err), C.ub4(pos+1),
			unsafe.Pointer(&buffer[0]), C.sb4(len(buffer)), C.SQLT_STR, nil, nil, nil, 0, nil, C.OCI_DEFAULT)
		if result != C.OCI_SUCCESS {
			return ociGetError(stmt.conn.err)
		}
	}
	return nil
}
