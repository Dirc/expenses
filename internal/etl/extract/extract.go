package extract

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dirc/expenses/internal/models"
)

func ReadCSV(filePath string) ([]models.Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // CSV is comma-separated

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
