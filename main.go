// Ref: https://universalglue.dev/posts/csv-to-sqlite/
package main

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var dbName string

func importCsv(csvName string, dbName string) {
    // 	var csvName string = "small-with-columns.csv"

	// Initialise DB
    // 	db := connectDB(dbName)
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatalf("ping failed: %s", err)
	}

	// Create table
	// Columns "bedrag" and "saldoNaBoeking" are defined as TEXT since they have a comma as seperator (a dot is needed for type REAL/float64).
	// We will convert it to a float64 when we start doing calculations
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS expenses (boekdatum TEXT, rekeningnummer TEXT, bedrag REAL, debetCredit TEXT, naamTegenrekening TEXT, tegenrekening TEXT, code TEXT, omschrijving TEXT, saldoNaBoeking REAL, transactionType TEXT)")
	if err != nil {
		log.Fatalf("prepare failed: %s", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("exec failed: %s", err)
	}

	// Read csv
	f, err := os.Open(csvName)
	if err != nil {
		log.Fatalf("open failed: %s", err)
	}
	r := csv.NewReader(f)
	// Read the header row.
	_, err = r.Read()
	if err != nil {
		log.Fatalf("missing header row(?): %s", err)
	}

	// Identify csv columns
	for {
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}

        boekdatum           := record[0]
        rekeningnummer      := record[1]
        bedrag              := record[2]
        debetCredit         := record[3]
        naamTegenrekening   := record[4]
        tegenrekening       := record[5]
        code                := record[6]
        omschrijving        := record[7]
        saldoNaBoeking      := record[8]
        //transactionType     := record[9] // extra column, not in csv

		// Map csv columns to DB columns
		stmt, err = db.Prepare("insert into expenses(boekdatum, rekeningnummer, bedrag, debetCredit, naamTegenrekening, tegenrekening, code, omschrijving, saldoNaBoeking, transactionType) values(?, ?, ?, ?, ?, ?, ?, ?, ?, 'unknown')")
		if err != nil {
			log.Fatalf("insert prepare failed: %s", err)
		}

		_, err = stmt.Exec(boekdatum, rekeningnummer, bedrag, debetCredit, naamTegenrekening, tegenrekening, code, omschrijving, saldoNaBoeking)
		if err != nil {
			log.Fatalf("insert failed(%s, %s, %s): %s", boekdatum, rekeningnummer, bedrag, err)
		}
	}
    fmt.Printf("CSV file %s imported successfully as %s!\n", csvName, dbName)
}

func cleanDB(dbName string) {
    err := os.Remove(dbName)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("File %s deleted successfully!\n", dbName)
}
func printTable(dbName string) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    rows, err := db.Query("SELECT boekdatum, rekeningnummer, bedrag, debetCredit, naamTegenrekening, tegenrekening, code, omschrijving, saldoNaBoeking, transactionType FROM expenses")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var boekdatum string
        var rekeningnummer string
        var bedrag string
        var debetCredit string
        var naamTegenrekening string
        var tegenrekening string
        var code string
        var omschrijving string
        var saldoNaBoeking string
        var transactionType string
        if err := rows.Scan(&boekdatum, &rekeningnummer, &bedrag, &debetCredit, &naamTegenrekening, &tegenrekening, &code, &omschrijving, &saldoNaBoeking, &transactionType); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("boekdatum: %s, rekeningnummer: %s, bedrag: %s, debetCredit: %s, naamTegenrekening: %s, tegenrekening: %s, code: %s, omschrijving: %s, saldoNaBoeking: %s, transactionType: %s\n", boekdatum, rekeningnummer, bedrag, debetCredit, naamTegenrekening, tegenrekening, code, omschrijving, saldoNaBoeking, transactionType)
    }

    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
}

func generateTransactionType(dbName string) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    transactionTypeMap := make(map[string][]string)
    transactionTypeMap["transactionType"] = []string{"Boodschappen"}
    transactionTypeMap["naamTegenrekening"] = []string{"PICNIC BY BUCKAROO"}
    transactionTypeMap["tegenrekening"] = []string{"tegenrekening"}
    transactionTypeMap["omschrijving"] = []string{"Bakkerij Neplenbroek%", "ALBERT HEIJN%"}

//     fmt.Println(m["key1"]) // prints ["Alice", "Bob", "Charlie"]
//     fmt.Println(m["key2"]) // prints ["Dave", "Eve", "Frank"]

    // General loop -> doesn't work!
//     for dbColumn, dbColumnValues := range transactionTypeMap {
//         // loop over all keys except the transactionType
//         if dbColumn == "transactionType" {
//             break
//         }
//         for _, value := range dbColumnValues {
//             _, err = db.Exec("UPDATE expenses SET transactionType = 'Boodschappen' WHERE naamTegenrekening LIKE (?)", value)
//             if err != nil {
//                 log.Fatal(err)
//             }
//         }
//     }

    // naamTegenrekening
    for _, value := range transactionTypeMap["naamTegenrekening"] {
        _, err = db.Exec("UPDATE expenses SET transactionType = 'Boodschappen' WHERE naamTegenrekening LIKE (?)", value)
        if err != nil {
            log.Fatal(err)
        }
    }
    // omschrijving
    for _, value := range transactionTypeMap["omschrijving"] {
        _, err = db.Exec("UPDATE expenses SET transactionType = 'Boodschappen' WHERE omschrijving LIKE (?)", value)
        if err != nil {
            log.Fatal(err)
        }
    }

    fmt.Println("Values inserted successfully!")
}


// ----- refactor ------

func executeSql(sqlQuery string) {
	var dbName string = "expenses.db"

	// Initialise DB
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatalf("ping failed: %s", err)
	}

	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		log.Fatalf("prepare failed: %s", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("exec failed: %s", err)
	}
}

// func returnSqlResults(sqlQuery string) {
func returnSqlResults() {
	var dbName string = "expenses.db"

	// Initialise DB
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatalf("ping failed: %s", err)
	}

// 	stmt, err := db.Prepare(sqlQuery)
// 	if err != nil {
// 		log.Fatalf("prepare failed: %s", err)
// 	}
//
// 	response, err := stmt.Query()
// 	if err != nil {
// 		log.Fatalf("exec failed: %s", err)
// 	}

	rows, err := db.Query("select 'Naam tegenrekening' from expenses")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    // Loop through the first result set.
    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err != nil {
            log.Fatal(err)
        }
        fmt.Println("Naam: %s\n", name)
    }

    // Advance to next result set.
//     rows.NextResultSet()

    // Loop through the second result set.
//     for rows.Next() {
//         fmt.Println(rows)
//     }

    // Check for any error in either result set.
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
// 	fmt.Println(response)
}

type transactionTypeTotal struct {
  transactionType string
  value float32
}

// func (transactionTypeTotal *transactionTypeTotal) getTotal(transactionType string) {
//     value =
// }

func main() {
  dbName = "expenses.db"
  cleanDB(dbName)
  importCsv("small-with-columns.csv", dbName)
  generateTransactionType(dbName)
  printTable(dbName)
//   executeSql("select * from expenses")
//   returnSql("select bedrag from expenses WHERE omschrijving LIKE '%PICNIC%'")
//   returnSql("SELECT SUM(bedrag) FROM expenses WHERE debetCredit = 'Debet'")
//   getDebet()
//   returnSqlResults()
}