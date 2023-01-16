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

func createTable(db *sql.DB) error {
	// Create table
	// Columns "bedrag" and "saldoNaBoeking" are defined as TEXT since they have a comma as seperator (a dot is needed for type REAL/float64).
	// We will convert it to a float64 when we start doing calculations
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS expenses (
	    boekdatum TEXT,
	    rekeningnummer TEXT,
	    bedrag REAL,
	    debetCredit TEXT,
	    naamTegenrekening TEXT,
	    tegenrekening TEXT,
	    code TEXT,
	    omschrijving TEXT,
	    saldoNaBoeking REAL,
	    transactionType TEXT
	)`)
	if err != nil {
		log.Fatalf("prepare failed: %s", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("exec failed: %s", err)
	}
	return nil
}

func importCsv(csvName string, db *sql.DB) error {
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
		stmt, err := db.Prepare("insert into expenses(boekdatum, rekeningnummer, bedrag, debetCredit, naamTegenrekening, tegenrekening, code, omschrijving, saldoNaBoeking, transactionType) values(?, ?, ?, ?, ?, ?, ?, ?, ?, 'unknown')")
		if err != nil {
			log.Fatalf("insert prepare failed: %s", err)
		}

		_, err = stmt.Exec(boekdatum, rekeningnummer, bedrag, debetCredit, naamTegenrekening, tegenrekening, code, omschrijving, saldoNaBoeking)
		if err != nil {
			log.Fatalf("insert failed(%s, %s, %s): %s", boekdatum, rekeningnummer, bedrag, err)
		}
	}
    fmt.Printf("CSV file %s imported successfully!\n", csvName)
    return nil
}

func cleanDB(dbName string) {
    err := os.Remove(dbName)
    if err != nil {
        if os.IsNotExist(err) {
            fmt.Println("File does not exist")
        } else {
            panic(err)
        }
    }
    fmt.Printf("File %s deleted successfully!\n", dbName)
}

func printTable(db *sql.DB) error {
    rows, err := db.Query("SELECT boekdatum, rekeningnummer, bedrag, debetCredit, naamTegenrekening, tegenrekening, code, omschrijving, saldoNaBoeking, transactionType FROM expenses")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    // Newline before printing table
    fmt.Printf("\n")

    // Print table
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
        return fmt.Errorf("Error printing table: %v", err)
    }
    return nil
}
