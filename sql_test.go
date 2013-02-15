package goci_test

import (
	_ "github.com/egravert/goci"
	"database/sql"
	"os"
	"testing"
)

func TestCanQueryDb(t *testing.T) {
	dsn := os.Getenv("ORACLE_DSN")
	db, err := sql.Open("goci", dsn)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Query("SELECT 'goci' FROM dual")
	if err != nil {
		t.Fatal(err)
	}
  
}
