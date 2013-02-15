package goci

/*
#include <oci.h>
#include <stdlib.h>
#include <string.h>

#cgo pkg-config: oci8
*/
import "C"
import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"unsafe"
)

type drv struct{}

func init() {
	sql.Register("goci", &drv{})
}

func (d *drv) Open(dsn string) (driver.Conn, error) {
	conn := &connection{}

	// initialize the oci environment used for all other oci calls
	result := C.OCIEnvCreate((**C.OCIEnv)(unsafe.Pointer(&conn.env)), C.OCI_DEFAULT, nil, nil, nil, nil, 0, nil)
	if result != C.OCI_SUCCESS {
		return nil, errors.New("Failed: OCIEnvCreate()")
	}

	// error handle
	result = C.OCIHandleAlloc(conn.env, &conn.err, C.OCI_HTYPE_ERROR, 0, nil)
	if result != C.OCI_SUCCESS {
		return nil, errors.New("Failed: OCIHandleAlloc() - creating error handle")
	}

	// Log in the user
	err := conn.performLogon(dsn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func ociGetError(err unsafe.Pointer) error {
	var errcode C.sb4
	var errbuff [512]C.char
	C.OCIErrorGet(err, 1, nil, &errcode, (*C.OraText)(unsafe.Pointer(&errbuff[0])), 512, C.OCI_HTYPE_ERROR)
	s := C.GoString(&errbuff[0])
	return errors.New(s)
}
