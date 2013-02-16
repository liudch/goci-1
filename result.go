package goci

/*
#include <oci.h>
#include <stdlib.h>
#include <string.h>

#cgo pkg-config: oci8
*/
import "C"
import (
	"unsafe"
)

type result struct {
	stmt *statement
}

func (r *result) LastInsertId() (int64, error) {
	var t C.ub4
	if C.OCIAttrGet(r.stmt.handle, C.OCI_HTYPE_STMT, unsafe.Pointer(&t), nil, C.OCI_ATTR_ROWID, (*C.OCIError)(r.stmt.conn.err)) != C.OCI_SUCCESS {
		return 0, ociGetError(r.stmt.conn.err)
	}
	return int64(t), nil
}

func (r *result) RowsAffected() (int64, error) {
	var t C.ub4
	if C.OCIAttrGet(r.stmt.handle, C.OCI_HTYPE_STMT, unsafe.Pointer(&t), nil, C.OCI_ATTR_ROW_COUNT, (*C.OCIError)(r.stmt.conn.err)) != C.OCI_SUCCESS {
		return 0, ociGetError(r.stmt.conn.err)
	}
	return int64(t), nil
}
