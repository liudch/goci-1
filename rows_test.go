package goci_test

import (
	"database/sql"
	_ "github.com/egravert/goci"
	"os"
	"testing"
)

func TestColumnNames(t *testing.T) {
	dsn := os.Getenv("ORACLE_DSN")
	db, err := sql.Open("goci", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT 1 as field1, 2 as field2, 3 as field3 FROM dual")
	if err != nil {
		t.Fatal(err)
	}

	columns, _ := rows.Columns()
	if len(columns) != 3 {
		t.Fatal("Expected 3 columns, got ", len(columns))
	}
	if columns[0] != "FIELD1" {
		t.Fatal("expected field1, got ", columns[0])
	}
	if columns[1] != "FIELD2" {
		t.Fatal("expected field2, got ", columns[1])
	}
	if columns[2] != "FIELD3" {
		t.Fatal("expected field3, got ", columns[0])
	}

	defer rows.Close()
}
