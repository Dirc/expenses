package main

import (
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type TransactionType struct {
    Type            string
    SearchValues    map[string][]string
}

func generateTransactionType(db *sql.DB, transTyp TransactionType) error {

    // General sql update
    for column := range transTyp.SearchValues {
        for _, searchString := range transTyp.SearchValues[column] {
            fmt.Printf("Column: %s, Value: %s, TransactionType: %s\n", column, searchString, transTyp.Type)

            // Build query
            query := "UPDATE expenses SET transactionType = '" + transTyp.Type + "' WHERE " + column + " LIKE '%" + searchString + "%'"

            // Prepare the statement
            stmt, err := db.Prepare(query)
            if err != nil {
                panic(err)
            }
            defer stmt.Close()

            // Execute the query
            result, err := stmt.Exec()
            if err != nil {
                panic(err)
            }

            // Get the number of affected rows
            affected, err := result.RowsAffected()
            if err != nil {
                panic(err)
            }
            fmt.Printf("Number of affected rows: %d\n", affected)
        }
    }

    fmt.Println("Values inserted successfully!")
    return nil
}
