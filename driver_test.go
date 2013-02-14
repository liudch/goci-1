package goci

import (
	"os"
	"testing"
)

// test setting up the base driver struct
func TestOpenReturnsValidConnection(t *testing.T) {
	driver := &drv{}

	dsn := os.Getenv("ORACLE_DSN")
	if dsn == "" {
		t.Fatal("To run tests, set the ORACLE_DSN environment variable.")
	}

	db, err := driver.Open(dsn)
	if err != nil {
		t.Fatal(err)
	}

	if db == nil {
		t.Fatal("db undefined")
	}
}

// a simple ping to oracle
func TestCanConnect(t *testing.T) {
	driver := &drv{}
	dsn := os.Getenv("ORACLE_DSN")
	if dsn == "" {
		t.Fatal("To run tests, set the ORACLE_DSN environment variable.")
	}

	db, _ := driver.Open(dsn)
	if err := db.(*conn).ping(); err != nil {
		t.Fatal(err)
	}
}

// test dsn parsing
func TestCanParseSimpleDSN(t *testing.T) {
	user, pwd, host := parseDsn("scott/tiger")
	
	stringsEqual(t, "scott", user)
	stringsEqual(t, "tiger", pwd)
	stringsEqual(t, "", host)
}

func TestCanParseFullDSN(t *testing.T) {
	user, pwd, host := parseDsn("scott/tiger@mycompany.com:2343/CRP01")
	stringsEqual(t, "scott", user)
	stringsEqual(t, "tiger", pwd)
	stringsEqual(t, "mycompany.com:2343/CRP01", host)
}

// should an error return?
func TestInvalidDSN(t *testing.T) {
	user, pwd, host := parseDsn("")
	stringsEqual(t, "", user)
	stringsEqual(t, "", pwd)
	stringsEqual(t, "", host)
}

func stringsEqual(t *testing.T, expected string, actual string) {
	if expected != actual {
		t.Fatal("Expected: '" + expected + "' Actual: '" + actual + "'")
	}
}
