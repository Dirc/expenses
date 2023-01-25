package main

import (
    "testing"
    "database/sql"

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

func TestImportCsv(t *testing.T) {
    // create a temporary SQLite3 database
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("error creating temporary db: %v", err)
    }
    defer db.Close()

    // create a table with the same schema as the one you're going to insert the data into
    _, err = db.Exec("CREATE TABLE expenses (boekdatum text, rekeningnummer text, bedrag text, debetCredit text, naamTegenrekening text, tegenrekening text, code text, omschrijving text, saldoNaBoeking text, transactionType text);")
    if err != nil {
        t.Fatalf("error creating table: %v", err)
    }

    // Read the test CSV file and set it as the input for the importCsv function
    csvName := "./test/resources/dummy.csv"
    err = importCsv(csvName, db)
    if err != nil {
        t.Fatalf("error importing csv: %v", err)
    }

    // Compare the result of the function with the expected result
    rows, err := db.Query("SELECT COUNT(*) FROM expenses")
    if err != nil {
        t.Fatalf("error selecting data from table: %v", err)
    }
    defer rows.Close()

    var count int
    var expected int
    expected = 10
    for rows.Next() {
        if err := rows.Scan(&count); err != nil {
            t.Fatalf("error scanning rows: %v", err)
        }
    }
    if count != expected {
        t.Errorf("The number of inserted rows is incorrect, got %d, want %d", count, expected)
    }
}