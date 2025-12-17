// Package extract extracts data into memory.
package extract

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dirc/expenses/internal/models"
)

// ReadCSV read csv file containing transactions.
func ReadCSV(filePath string) ([]models.Transaction, error) {
	// if !isSafePath(filePath) {
	// 	return fmt.Errorf("invalid file path: %s", filePath)
	// }
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}()

	reader := csv.NewReader(file)
	reader.Comma = ','

	// Read and discard header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction

	for {
		record, err := reader.Read()
		if err != nil {
			break // EOF or error
		}

		// Parse Bedrag (replace comma with dot for float parsing)
		bedrag, err := strconv.ParseFloat(
			strings.ReplaceAll(record[2], ",", "."),
			64,
		)
		if err != nil {
			return nil, err
		}

		// Parse Boekdatum (DD-MM-YYYY)
		boekdatum, err := time.Parse("02-01-2006", record[0])
		if err != nil {
			return nil, err
		}

		// Parse SaldoNaBoeking (replace comma with dot for float parsing)
		saldo, err := strconv.ParseFloat(
			strings.ReplaceAll(record[8], ",", "."),
			64,
		)
		if err != nil {
			return nil, err
		}

		tx := models.Transaction{
			Boekdatum:         boekdatum,
			Rekeningnummer:    record[1],
			Bedrag:            bedrag,
			DebetCredit:       record[3],
			NaamTegenrekening: record[4],
			Tegenrekening:     record[5],
			Code:              record[6],
			Omschrijving:      record[7],
			SaldoNaBoeking:    saldo,
			TransactionType:   "",
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}
