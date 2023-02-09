package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type TransactionType struct {
	Type         string              `yaml:"type"`
	SearchValues map[string][]string `yaml:"SearchValues"`
}

type TransactionTypes struct {
	TransactionTypes []TransactionType `yaml:"TransactionTypes"`
}

func readTransactionTypesFromYaml(config string) (TransactionTypes, error) {
	file, err := os.Open(config)
	if err != nil {
		fmt.Printf("Error reading YAML file: %v", err)
		return TransactionTypes{}, err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading contents of YAML file: %v", err)
		return TransactionTypes{}, err
	}

	transactionTypes := TransactionTypes{}
	err = yaml.Unmarshal(fileBytes, &transactionTypes)
	if err != nil {
		fmt.Println(err)
		return TransactionTypes{}, err
	}

	return transactionTypes, nil
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

func generateAllTransactionTypes(configfile string, db *sql.DB) error {
	t, err := readTransactionTypesFromYaml(configfile)
	if err != nil {
		fmt.Println(err)
	}

	for _, transType := range t.TransactionTypes {
		fmt.Println(transType.Type)
		generateTransactionType(db, transType)
	}
	return nil
}

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

func PrintTableTransactionType(db *sql.DB, transactionType string) error {

	var query string
	var rows *sql.Rows
	var err error

	if transactionType == "" {
		query = "SELECT boekdatum, rekeningnummer, bedrag, debetCredit, naamTegenrekening, tegenrekening, code, omschrijving, saldoNaBoeking, transactionType FROM expenses;"
		rows, err = db.Query(query)
	} else {
		query = "SELECT boekdatum, rekeningnummer, bedrag, debetCredit, naamTegenrekening, tegenrekening, code, omschrijving, saldoNaBoeking, transactionType FROM expenses WHERE transactionType = ?;"
		rows, err = db.Query(query, transactionType)
	}

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

// For clearity, not really needed
func PrintTable(db *sql.DB) error {
	return PrintTableTransactionType(db, "")
}
