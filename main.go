// Ref: https://universalglue.dev/posts/csv-to-sqlite/
package main

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
// 	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

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

	// Create table
	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		log.Fatalf("prepare failed: %s", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("exec failed: %s", err)
	}
}

func importCsv(csvName string) {
	var dbName string = "expenses.db"
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
	stmt, err := db.Prepare("create table if not exists expenses (boekdatum text, rekeningnummer text, bedrag real, debetCredit text, naamTegenrekening text, tegenrekening text, code text, omschrijving text, saldoNaBoeking text, transactionType text)")
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
}

func main() {
  importCsv("small-with-columns.csv")
  executeSql("select * from expenses")
//   getDebet()
}