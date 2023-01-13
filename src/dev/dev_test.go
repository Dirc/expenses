package main

import (
    "database/sql"
    "testing"

    _ "github.com/mattn/go-sqlite3"
)

func TestCreateTable(t *testing.T) {

    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Errorf("Error opening database: %v", err)
    }
    defer db.Close()

    // Create the table
    createTable(db)

    // Check that the table exists
    rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='users'")
    if err != nil {
        t.Errorf("Error querying table: %v", err)
    }
    if !rows.Next() {
        t.Errorf("Table 'users' not found")
    }
}

func TestFoo(t *testing.T) {
    // Test that the foo function returns the expected result
    result := foo()
    if result != 42 {
        t.Errorf("Unexpected result: %d", result)
    }
}