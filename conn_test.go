package goci

import (
	"os"
	"testing"
)

func TestCanConnect(t *testing.T) {
	driver := &drv{}
	db, err := driver.Open(os.Getenv("ORACLE_DSN"))
	if err != nil {
		t.Fatal(err)
	}

	if db == nil {
		t.Fatal("db undefined")
	}

	if err = db.(*conn).ping(); err != nil {
		t.Fatal(err)
	}

}
