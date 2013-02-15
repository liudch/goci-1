package goci_test

import (
	"database/sql"
	_ "github.com/egravert/goci"
	"os"
	"testing"
)

func TestStatements(t *testing.T) {
	dsn := os.Getenv("ORACLE_DSN")
	db, err := sql.Open("goci", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT cast(1 as int) as foo, 3.14 as bar, 'goci' as baz FROM dual")
	if err != nil {
		t.Fatal(err)
	}

	_, err = stmt.Query()
	if err != nil {
		t.Fatal(err)
	}
  //defer rows.Close()
}
