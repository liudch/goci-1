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
	"strings"
	"unsafe"
)

type conn struct {
	env_handle    unsafe.Pointer
	error_handle  unsafe.Pointer
	server_handle unsafe.Pointer
}

func (cn *conn) performLogon(dsn string) error {

	user, pwd, host := parseDsn(dsn)
	c_user := C.CString(user)
	defer C.free(unsafe.Pointer(c_user))
	c_pwd := C.CString(pwd)
	defer C.free(unsafe.Pointer(c_pwd))
	c_host := C.CString(host)
	defer C.free(unsafe.Pointer(c_host))

	result := C.OCILogon2((*C.OCIEnv)(unsafe.Pointer(cn.env_handle)),
		(*C.OCIError)(cn.error_handle),
		(**C.OCIServer)(unsafe.Pointer(&cn.server_handle)),
		(*C.OraText)(unsafe.Pointer(c_user)),
		C.ub4(C.strlen(c_user)),
		(*C.OraText)(unsafe.Pointer(c_pwd)),
		C.ub4(C.strlen(c_pwd)),
		(*C.OraText)(unsafe.Pointer(c_host)),
		C.ub4(C.strlen(c_host)),
		C.OCI_LOGON2_STMTCACHE)
	if result != C.OCI_SUCCESS {
		return ociGetError(cn.error_handle)
	}
	return nil
}

func (cn *conn) Begin() (driver.Tx, error) {
	return nil, nil
}

func (cn *conn) Prepare(query string) (driver.Stmt, error) {
	pquery := C.CString(query)
	defer C.free(unsafe.Pointer(pquery))
	var stmt_handle unsafe.Pointer

	if C.OCIHandleAlloc(cn.env_handle, &stmt_handle, C.OCI_HTYPE_STMT, 0, nil) != C.OCI_SUCCESS {
		return nil, ociGetError(cn.error_handle)
	}
	return &statement{handle: stmt_handle}, nil
}

func (cn *conn) Close() error {
	return nil
}

// Makes a lightweight call to the server. A successful result indicates the server is active.  A block indicates the connection may be in use by
// another thread. A failure indicates a communication error.
func (cn *conn) ping() error {
	if C.OCIPing((*C.OCIServer)(cn.server_handle), (*C.OCIError)(cn.error_handle), C.OCI_DEFAULT) != C.OCI_SUCCESS {
		return ociGetError(cn.error_handle)
	}
	return nil
}

// expect the dsn in the format of: user/pwd@host:port/SID
func parseDsn(dsn string) (user string, pwd string, host string) {
	tokens := strings.SplitN(dsn, "@", 2)
	if len(tokens) > 1 {
		host = tokens[1]
	}
	userpass := strings.SplitN(tokens[0], "/", 2)
	if len(userpass) > 1 {
		pwd = userpass[1]
	}
	user = userpass[0]
	return user, pwd, host
}
