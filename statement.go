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
	"log"
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
	var num C.int
	if r := C.OCIAttrGet(stmt.handle, C.OCI_HTYPE_STMT, unsafe.Pointer(&num), nil, C.OCI_ATTR_BIND_COUNT, (*C.OCIError)(stmt.conn.err)); r != C.OCI_SUCCESS {
		log.Println(ociGetError(stmt.conn.err))
	}
	return int(num)
}

// Exec executes a query that doesn't return rows, such
// as an INSERT or UPDATE.
func (stmt *statement) Exec(v []driver.Value) (driver.Result, error) {
	return nil, nil
}

// Exec executes a query that may return rows, such as SELECT.
func (stmt *statement) Query(v []driver.Value) (rows driver.Rows, err error) {
	if err = stmt.bind(v); err != nil {
		return
	}

	// determine the type of statement.  For select statements iter is set to zero, for other statements, only execute once.
	// * NOTE * Should an error be returned if the statement is a non-select type? 
	var stmt_type C.int
	if C.OCIAttrGet(stmt.handle, C.OCI_HTYPE_STMT, unsafe.Pointer(&stmt_type), nil, C.OCI_ATTR_STMT_TYPE, (*C.OCIError)(stmt.conn.err)) != C.OCI_SUCCESS {
		log.Println(ociGetError(stmt.conn.err))
	}
	iter := C.ub4(1)
	if stmt_type == C.OCI_STMT_SELECT {
		iter = 0
	}
	// set the row prefetch.  Only one extra row per fetch will be returned unless this is set.
	prefetchSize := C.ub4(100)
	if C.OCIAttrSet(stmt.handle, C.OCI_HTYPE_STMT, unsafe.Pointer(&prefetchSize), 0, C.OCI_ATTR_PREFETCH_ROWS, (*C.OCIError)(stmt.conn.err)) != C.OCI_SUCCESS {
		log.Println(ociGetError(stmt.conn.err))
	}

	// execute the statement
	if C.OCIStmtExecute((*C.OCIServer)(stmt.conn.svr), (*C.OCIStmt)(stmt.handle), (*C.OCIError)(stmt.conn.err), iter, 0, nil, nil, C.OCI_DEFAULT) != C.OCI_SUCCESS {
		err = ociGetError(stmt.conn.err)
		log.Println(err)
		return
	}

	// find out how many output columns there are
	var cols C.ub2
	if C.OCIAttrGet(stmt.handle, C.OCI_HTYPE_STMT, unsafe.Pointer(&cols), nil, C.OCI_ATTR_PARAM_COUNT, (*C.OCIError)(stmt.conn.err)) != C.OCI_SUCCESS {
		err = ociGetError(stmt.conn.err)
		log.Println(err)
		return
	}

	// build column meta-data
	columns := make([]column, int(cols))
	for pos, column := range columns {
		var param unsafe.Pointer
		var colType C.ub2
		var colSize C.ub4
		var colName *C.char
		var nameSize C.ub4

		if C.OCIParamGet(stmt.handle, C.OCI_HTYPE_STMT, (*C.OCIError)(stmt.conn.err), (*unsafe.Pointer)(unsafe.Pointer(&param)), C.ub4(pos+1)) != C.OCI_SUCCESS {
			err = ociGetError(stmt.conn.err)
			log.Println(err)
			return
		}

		C.OCIAttrGet(param, C.OCI_DTYPE_PARAM, unsafe.Pointer(&colType), nil, C.OCI_ATTR_DATA_TYPE, (*C.OCIError)(stmt.conn.err))
		C.OCIAttrGet(param, C.OCI_DTYPE_PARAM, unsafe.Pointer(&colName), &nameSize, C.OCI_ATTR_NAME, (*C.OCIError)(stmt.conn.err))
		C.OCIAttrGet(param, C.OCI_DTYPE_PARAM, unsafe.Pointer(&colSize), nil, C.OCI_ATTR_DATA_SIZE, (*C.OCIError)(stmt.conn.err))

    column.kind = int(colType)
    column.size = int(colSize)
    column.name = C.GoStringN(colName, (C.int)(nameSize))
    fmt.Println(column)
    //column.raw = make([]byte, int(colSize))
	}
	return
}

type column struct {
	name string
	kind int
	size int
	raw []byte
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
