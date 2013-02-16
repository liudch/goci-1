package goci_test

import (
	"database/sql"
	_ "github.com/egravert/goci"
	"os"
	"testing"
  "fmt"
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

	rows, err := stmt.Query()
	if err != nil {
		t.Fatal(err)
	}

	var item1, item2, item3 string
  fmt.Println("here")
	for rows.Next() {
		rows.Scan(&item1, &item2, &item3)
		fmt.Println(item1, item2, item3)
	}

	defer rows.Close()
}
