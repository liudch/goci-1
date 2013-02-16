package goci

/*
#include <oci.h>
#include <stdlib.h>
#include <string.h>

#cgo pkg-config: oci8
*/
import "C"

type transaction struct {
	conn *connection
}

func (tx *transaction) Commit() error {
	if err := tx.conn.exec("COMMIT"); err != nil {
		return err
	}
	return nil
}

func (tx *transaction) Rollback() error {
	if err := tx.conn.exec("ROLLBACK"); err != nil {
		return err
	}
	return nil
}
