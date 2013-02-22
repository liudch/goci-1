// Package oci provides a light wrapper over the OCI native methdods.  All methods check for errors and return an errors object if necessary.
// The returned structures provide type safety for the unsafe pointers required for native OCI calls.
package oci

/*
#include <oci.h>
#include <stdlib.h>
#include <string.h>

#cgo pkg-config: oci8
*/
import "C"
import (
	"errors"
	"strconv"
	"unsafe"
)

type ociHandle struct {
	handle unsafe.Pointer
	htype  C.ub4
}

// The environment under which all other structures reside.
type EnvHandle ociHandle

// An error 
type ErrHandle ociHandle

// A pointer to the server a login has been established with.
type SvrHandle ociHandle

// A handle to a prepared statement.  The statement can be used to bind parameters and be executed.
type StmtHandle ociHandle

// CreateEnvironment initializes the environment that OCI will work under. The environment is configured for use in an threaded environment.
func CreateEnvironment() (env EnvHandle, err error) {
	result := C.OCIEnvCreate((**C.OCIEnv)(unsafe.Pointer(&env.handle)), C.OCI_THREADED, nil, nil, nil, nil, 0, nil)
	if result != C.OCI_SUCCESS {
		err = errors.New("Failed: OCIEnvCreate()")
	} else {
		env.htype = C.OCI_HTYPE_ENV
	}
	return
}

// Allocates an error handle used for most other method calls.
func AllocateErrorHandle(env EnvHandle) (ociErr ErrHandle, err error) {
	ociErr.htype = C.OCI_HTYPE_ERROR
	result := C.OCIHandleAlloc(env.handle, &ociErr.handle, ociErr.htype, 0, nil)
	if result != C.OCI_SUCCESS {
		err = errors.New("Failed: OCIHandleAlloc() - creating error handle. Error code: " + strconv.Itoa(int(result)))
	}
	return
}

// Releases the error handle
func FreeErrorHandle(handle ErrHandle) (err error) {
	result := C.OCIHandleFree(handle.handle, handle.htype)
	if result != C.OCI_SUCCESS {
		err = errors.New("Failed: OCIHandleFree() - freeing error handle. Error code: " + strconv.Itoa(int(result)))
	}
	return
}

// Allocate a new statement handle of type OCI_HTYPE_STMT
func AllocateStatementHandle(env EnvHandle) (stmt StmtHandle, err error) {
	stmt.htype = C.OCI_HTYPE_STMT
	result := C.OCIHandleAlloc(env.handle, &stmt.handle, stmt.htype, 0, nil)
	if result != C.OCI_SUCCESS {
		err = errors.New("Failed: OCIHandleAlloc() - creating error handle. Error code: " + strconv.Itoa(int(result)))
	}
	return
}

// Releases the statement handle
func FreeStatementHandle(handle StmtHandle) (err error) {
	result := C.OCIHandleFree(handle.handle, handle.htype)
	if result != C.OCI_SUCCESS {
		err = errors.New("Failed: OCIHandleFree() - freeing statement handle. Error code: " + strconv.Itoa(int(result)))
	}
	return
}

// Performs a basic login to oracle.  The host may be a TSN or in the format of host:port/SID
func BasicLogin(env EnvHandle, user, pwd, host string) (svr SvrHandle, err error) {

	ociErr, err := AllocateErrorHandle(env)
	if err != nil {
		return
	}
	defer FreeErrorHandle(ociErr)

	puser := C.CString(user)
	defer C.free(unsafe.Pointer(puser))
	ppwd := C.CString(pwd)
	defer C.free(unsafe.Pointer(ppwd))
	phost := C.CString(host)
	defer C.free(unsafe.Pointer(phost))

	result := C.OCILogon2((*C.OCIEnv)(env.handle),
		(*C.OCIError)(ociErr.handle),
		(**C.OCIServer)(unsafe.Pointer(&svr.handle)),
		(*C.OraText)(unsafe.Pointer(puser)),
		C.ub4(C.strlen(puser)),
		(*C.OraText)(unsafe.Pointer(ppwd)),
		C.ub4(C.strlen(ppwd)),
		(*C.OraText)(unsafe.Pointer(phost)),
		C.ub4(C.strlen(phost)),
		C.OCI_LOGON2_STMTCACHE)
	if result != C.OCI_SUCCESS {
		err = ociGetError(ociErr)
	}
	return
}

// Creates a statement handle for the passed in statement.  
func Prepare(env EnvHandle, query string) (stmt StmtHandle, err error) {
	var ociErr ErrHandle

	pquery := C.CString(query)
	defer C.free(unsafe.Pointer(pquery))

	if ociErr, err = AllocateErrorHandle(env); err != nil {
		return
	}
	defer FreeErrorHandle(ociErr)

	if stmt, err = AllocateStatementHandle(env); err != nil {
		return
	}
	defer FreeErrorHandle(ociErr)

	result := C.OCIStmtPrepare((*C.OCIStmt)(stmt.handle), (*C.OCIError)(ociErr.handle), (*C.OraText)(unsafe.Pointer(pquery)), C.ub4(C.strlen(pquery)), C.ub4(C.OCI_NTV_SYNTAX), C.ub4(C.OCI_DEFAULT))
	if result != C.OCI_SUCCESS {
		err = ociGetError(ociErr)
	}
	return
}

// Makes a lightweight call to the server. A successful result indicates the server is active.  A block indicates the connection may be in use by
// another thread. A failure indicates a communication error.
func Ping(env EnvHandle, svr SvrHandle) error {
	ociErr, err := AllocateErrorHandle(env)
	if err != nil {
		return err
	}
	defer FreeErrorHandle(ociErr)
	if C.OCIPing((*C.OCIServer)(svr.handle), (*C.OCIError)(ociErr.handle), C.OCI_DEFAULT) != C.OCI_SUCCESS {
		return ociGetError(ociErr)
	}
	return nil
}

// Uses the error handle to return the textual error message returned from oracle.
func ociGetError(err ErrHandle) error {
	var errcode C.sb4
	var errbuff [512]C.char
	C.OCIErrorGet(err.handle, 1, nil, &errcode, (*C.OraText)(unsafe.Pointer(&errbuff[0])), 512, C.OCI_HTYPE_ERROR)
	s := C.GoString(&errbuff[0])
	return errors.New(s)
}
