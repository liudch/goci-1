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

type connection struct {
	env unsafe.Pointer
	err unsafe.Pointer
	svr unsafe.Pointer
}

func (conn *connection) performLogon(dsn string) error {

	user, pwd, host := parseDsn(dsn)
	puser := C.CString(user)
	defer C.free(unsafe.Pointer(puser))
	ppwd := C.CString(pwd)
	defer C.free(unsafe.Pointer(ppwd))
	phost := C.CString(host)
	defer C.free(unsafe.Pointer(phost))

	result := C.OCILogon2((*C.OCIEnv)(unsafe.Pointer(conn.env)),
		(*C.OCIError)(conn.err),
		(**C.OCIServer)(unsafe.Pointer(&conn.svr)),
		(*C.OraText)(unsafe.Pointer(puser)),
		C.ub4(C.strlen(puser)),
		(*C.OraText)(unsafe.Pointer(ppwd)),
		C.ub4(C.strlen(ppwd)),
		(*C.OraText)(unsafe.Pointer(phost)),
		C.ub4(C.strlen(phost)),
		C.OCI_LOGON2_STMTCACHE)
	if result != C.OCI_SUCCESS {
		return ociGetError(conn.err)
	}
	return nil
}

func (conn *connection) exec(cmd string) error {
  stmt, err := conn.Prepare(cmd)
  if err == nil {
    defer stmt.Close()
    _, err = stmt.Exec(nil)
  }
  return err
}

func (conn *connection) Begin() (driver.Tx, error) {
  if err := conn.exec("BEGIN"); err != nil {
    return nil, err
  }
  return &transaction{conn}, nil
}

func (conn *connection) Prepare(query string) (driver.Stmt, error) {
	pquery := C.CString(query)
	defer C.free(unsafe.Pointer(pquery))
	var stmt unsafe.Pointer

	if C.OCIHandleAlloc(conn.env, &stmt, C.OCI_HTYPE_STMT, 0, nil) != C.OCI_SUCCESS {
		return nil, ociGetError(conn.err)
	}
	result := C.OCIStmtPrepare((*C.OCIStmt)(stmt), (*C.OCIError)(conn.err),
		(*C.OraText)(unsafe.Pointer(pquery)), C.ub4(C.strlen(pquery)),
		C.ub4(C.OCI_NTV_SYNTAX), C.ub4(C.OCI_DEFAULT))
	if result != C.OCI_SUCCESS {
		return nil, ociGetError(conn.err)
	}
	return &statement{handle: stmt, conn: conn}, nil
}

func (conn *connection) Close() error {
	return nil
}

// Makes a lightweight call to the server. A successful result indicates the server is active.  A block indicates the connection may be in use by
// another thread. A failure indicates a communication error.
func (conn *connection) ping() error {
	if C.OCIPing((*C.OCIServer)(conn.svr), (*C.OCIError)(conn.err), C.OCI_DEFAULT) != C.OCI_SUCCESS {
		return ociGetError(conn.err)
	}
	return nil
}

// expect the dsn in the format of: user/pwd@host:port/SID
func parseDsn(dsn string) (user, pwd, host string) {
	tokens := strings.SplitN(dsn, "@", 2)
	if len(tokens) > 1 {
		host = tokens[1]
	}
	userpass := strings.SplitN(tokens[0], "/", 2)
	if len(userpass) > 1 {
		pwd = userpass[1]
	}
	user = userpass[0]
	return
}
