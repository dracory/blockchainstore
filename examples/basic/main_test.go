package main

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestBasicExample(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	err = runBasicExample(db)
	if err != nil {
		t.Fatalf("runBasicExample failed: %v", err)
	}
}
