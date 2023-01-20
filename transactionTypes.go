package main

import (
	"fmt"
	"log"
	"strconv"
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

// Take print table as example!

func transactionTypesTotal(db *sql.DB) error {
    query := "SELECT transactionType, SUM(bedrag) as sum_score FROM expenses GROUP BY transactionType;"

    // Find the maximum length of each column
    rows, err := db.Query(query)
    if err != nil {
        fmt.Println("Error querying the database:", err)
    }
    defer rows.Close()

    // Max length of each database item
    maxLengths := make([]int, 2)
    for rows.Next() {
        var transactionType string
        var sum int
        err = rows.Scan(&transactionType, &sum)
        if err != nil {
            fmt.Println("Error reading rows:", err)
        }

        if len(transactionType) > maxLengths[0] {
            maxLengths[0] = len(transactionType)
        }
        sumStr := strconv.Itoa(sum)
        if len(sumStr) > maxLengths[1] {
            maxLengths[1] = len(sumStr)
        }
    }

    // Add length of table heading
    head1 := "transactionType"
    head2 := "sum"
    if len(head1) > maxLengths[0] {
        maxLengths[0] = len(head1)
    }
    if len(head2) > maxLengths[1] {
        maxLengths[1] = len(head2)
    }

    // Print table using maxLengths
    rows, err = db.Query(query)
    if err != nil {
        fmt.Println("Error querying the database:", err)
    }
    defer rows.Close()

    // Print table headers with the same left padding
    fmt.Printf("%-*s | %-*s\n", maxLengths[0], "transactionType", maxLengths[1], "sum")

    // Print rows with padding
    for rows.Next() {
        var transactionType string
        var sum int
        err = rows.Scan(&transactionType, &sum)
        if err := rows.Err(); err != nil {
            log.Fatal(err)
            return fmt.Errorf("Error printing table: %v", err)
        }
        // Print each cell with the same left padding
        fmt.Printf("%-*s | %-*d\n", maxLengths[0], transactionType, maxLengths[1], sum)
    }

    return nil
}