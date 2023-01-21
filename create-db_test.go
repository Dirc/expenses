package main

import (
    "database/sql"
    "testing"

    _ "github.com/mattn/go-sqlite3"
)

// Helper func
func contains(lst []string, str string) bool {
    for _, a := range lst {
        if a == str {
            return true
        }
    }
    return false
}

func TestCreateTable(t *testing.T) {
    expectedColumns := []string{"boekdatum", "rekeningnummer", "bedrag", "debetCredit", "naamTegenrekening", "tegenrekening", "code", "omschrijving", "saldoNaBoeking", "transactionType"}

    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Errorf("Error opening database: %v", err)
    }
    defer db.Close()

    // Create the table
    createTable(db)

    rows, err := db.Query("PRAGMA table_info(expenses)")
    if err != nil {
        t.Errorf("Error querying the database: %v", err)
    }
    defer rows.Close()

    actualColumns := []string{}
    for rows.Next() {
        var cid int
        var name, typ string
        var notnull, dfltVal, pk *int
        err = rows.Scan(&cid, &name, &typ, &notnull, &dfltVal, &pk)
        if err != nil {
            t.Errorf("Error scanning rows: %v", err)
        }
        actualColumns = append(actualColumns, name)
    }

    for _, column := range expectedColumns {
        if !contains(actualColumns, column) {
            t.Errorf("Table does not contain expected column %s", column)
        }
    }
}
