package main

import (
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
    // init db
    dbName := "expenses.db"
    fmt.Println(dbName)

    cleanDB(dbName)

    db, err := sql.Open("sqlite3", dbName)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    err = createTable(db)
    if err != nil {
        panic(err)
    }

    err = importCsv("small-with-columns.csv", db)
    if err != nil {
        panic(err)
    }

    // set transaction types
    boodschappen := TransactionType{
        Type: "boodschappen",
        SearchValues: map[string][]string{
            "naamTegenrekening": []string{"PICNIC BY BUCKAROO"},
            "tegenrekening": []string{"nog geen tegenrekening"},
            "omschrijving": []string{"Bakkerij Neplenbroek%", "ALBERT HEIJN%"},
        },
    }
    err = generateTransactionType(db, boodschappen)
    if err != nil {
        panic(err)
    }
    auto := TransactionType{
        Type: "auto",
        SearchValues: map[string][]string{
            "naamTegenrekening": []string{"none"},
            "tegenrekening": []string{"nog geen tegenrekening"},
            "omschrijving": []string{"Parkmobile Benelux BV%", "CCV*Auto-Veer%"},
        },
    }
    err = generateTransactionType(db, auto)
    if err != nil {
        panic(err)
    }
    spaarrek := TransactionType{
        Type: "spaarrek",
        SearchValues: map[string][]string{
            "naamTegenrekening": []string{"none"},
            "tegenrekening": []string{"TRIONL2U NL61TRIO2024035957"},
            "omschrijving": []string{"none"},
        },
    }
    err = generateTransactionType(db, spaarrek)
    if err != nil {
        panic(err)
    }
    gwl := TransactionType{
        Type: "gwl",
        SearchValues: map[string][]string{
            "naamTegenrekening": []string{"ENECO SERVICES"},
            "tegenrekening": []string{"none"},
            "omschrijving": []string{"none"},
        },
    }
    err = generateTransactionType(db, gwl)
    if err != nil {
        panic(err)
    }

    // Print status
//     err = printTable(db)
//     if err != nil {
//         panic(err)
//     }

    err = transactionTypesTotal(db)
    if err != nil {
        panic(err)
    }

    err = PrintTableTransactionType(db, "unknown")
    if err != nil {
        panic(err)
    }

    err = PrintTable(db)
    if err != nil {
        panic(err)
    }
}