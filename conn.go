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
	"strings"
	"unsafe"
)

type drv struct{}
type conn struct {
	env_handle    unsafe.Pointer
	error_handle  unsafe.Pointer
	server_handle unsafe.Pointer
}

func init() {
	sql.Register("goci", &drv{})
}

func (d *drv) Open(dsn string) (driver.Conn, error) {
	cn := &conn{}

	// initialize the oci environment used for all other oci calls
	result := C.OCIEnvCreate((**C.OCIEnv)(unsafe.Pointer(&cn.env_handle)), C.OCI_DEFAULT, nil, nil, nil, nil, 0, nil)
	if result != C.OCI_SUCCESS {
		return nil, errors.New("Failed: OCIEnvCreate()")
	}

	// error handle
	result = C.OCIHandleAlloc(cn.env_handle, &cn.error_handle, C.OCI_HTYPE_ERROR, 0, nil)
	if result != C.OCI_SUCCESS {
		return nil, errors.New("Failed: OCIHandleAlloc() - creating error handle")
	}

	// Log in the user
	err := cn.performLogon(dsn)
	if err != nil {
		return nil, err
	}
	return cn, nil
}

// expect the dsn in the format of: user/pwd@host:port/SID
func (cn *conn) performLogon(dsn string) error {
	tokens := strings.SplitN(dsn, "@", 2)
	userpass := strings.SplitN(tokens[0], "/", 2)

	var host *C.char
	hostlen := C.size_t(0)
	if len(tokens) > 1 {
		host = C.CString(tokens[1])
		defer C.free(unsafe.Pointer(host))
		hostlen = C.strlen(host)
	}
	user := C.CString(userpass[0])
	defer C.free(unsafe.Pointer(user))
	pass := C.CString(userpass[1])
	defer C.free(unsafe.Pointer(pass))

	result := C.OCILogon2((*C.OCIEnv)(unsafe.Pointer(cn.env_handle)),
		(*C.OCIError)(cn.error_handle),
		(**C.OCIServer)(unsafe.Pointer(&cn.server_handle)),
		(*C.OraText)(unsafe.Pointer(user)),
		C.ub4(C.strlen(user)),
		(*C.OraText)(unsafe.Pointer(pass)),
		C.ub4(C.strlen(pass)),
		(*C.OraText)(unsafe.Pointer(host)),
		C.ub4(hostlen),
		C.OCI_LOGON2_STMTCACHE)
	if result != C.OCI_SUCCESS {
		return ociGetError(cn.error_handle)
	}
	return nil
}

func (cn *conn) Begin() (driver.Tx, error) {
	return nil, nil
}

func (cn *conn) Close() error {
	return nil
}

func (cn *conn) Prepare(query string) (driver.Stmt, error) {
	return nil, nil
}

func (cn *conn) ping() error {
	if C.OCIPing((*C.OCIServer)(cn.server_handle), (*C.OCIError)(cn.error_handle), C.OCI_DEFAULT) != C.OCI_SUCCESS {
		return ociGetError(cn.error_handle)
	}
	return nil
}

func ociGetError(err unsafe.Pointer) error {
	var errcode C.sb4
	var errbuff [512]C.char
	C.OCIErrorGet(err, 1, nil, &errcode, (*C.OraText)(unsafe.Pointer(&errbuff[0])), 512, C.OCI_HTYPE_ERROR)
	s := C.GoString(&errbuff[0])
	return errors.New(s)
}
