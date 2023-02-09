package main

import (
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := Init("test/resources/dummy.csv")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = generateAllTransactionTypes("transactiontypes.yaml", db)
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
