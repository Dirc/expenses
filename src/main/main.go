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
    transTypBoodschappen := TransactionType{
        Type: "boodschappen",
        SearchValues: map[string][]string{
            "naamTegenrekening": []string{"PICNIC BY BUCKAROO"},
            "tegenrekening": []string{"nog geen tegenrekening"},
            "omschrijving": []string{"Bakkerij Neplenbroek%", "ALBERT HEIJN%"},
        },
    }

    err = generateTransactionType(db, transTypBoodschappen)
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
}