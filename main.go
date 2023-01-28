package main

import (
	"database/sql"
	"fmt"
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

	err = importCsv("test/resources/dummy.csv", db)
	if err != nil {
		panic(err)
	}

	// set transaction types
	boodschappen := TransactionType{
		Type: "boodschappen",
		SearchValues: map[string][]string{
			"naamTegenrekening": {"PICNIC BY BUCKAROO"},
			"tegenrekening":     {"nog geen tegenrekening"},
			"omschrijving":      {"Bakkerij Neplenbroek%", "ALBERT HEIJN%"},
		},
	}
	err = generateTransactionType(db, boodschappen)
	if err != nil {
		panic(err)
	}

	auto := TransactionType{
		Type: "auto",
		SearchValues: map[string][]string{
			"naamTegenrekening": {"none"},
			"tegenrekening":     {"nog geen tegenrekening"},
			"omschrijving":      {"Parkmobile Benelux BV%", "CCV*Auto-Veer%"},
		},
	}
	err = generateTransactionType(db, auto)
	if err != nil {
		panic(err)
	}

	spaarrek := TransactionType{
		Type: "spaarrek",
		SearchValues: map[string][]string{
			"naamTegenrekening": {"Dirc"},
			"tegenrekening":     {""},
			"omschrijving":      {"none"},
		},
	}
	err = generateTransactionType(db, spaarrek)
	if err != nil {
		panic(err)
	}

	gwl := TransactionType{
		Type: "gwl",
		SearchValues: map[string][]string{
			"naamTegenrekening": {"ENECO SERVICES"},
			"tegenrekening":     {"none"},
			"omschrijving":      {"none"},
		},
	}
	err = generateTransactionType(db, gwl)
	if err != nil {
		panic(err)
	}

	overig := TransactionType{
		Type: "overig",
		SearchValues: map[string][]string{
			"naamTegenrekening": {"M Grupetto"},
			"tegenrekening":     {"none"},
			"omschrijving":      {"none"},
		},
	}
	err = generateTransactionType(db, overig)
	if err != nil {
		panic(err)
	}

	err = transactionTypesTotal(db)
	if err != nil {
		panic(err)
	}

	err = PrintTableTransactionType(db, "unknown")
	if err != nil {
		panic(err)
	}

	//     err = PrintTable(db)
	//     if err != nil {
	//         panic(err)
	//     }
}
