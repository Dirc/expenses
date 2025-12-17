// Package transform enriches data.
package transform

import (
	"os"
	"strings"

	"github.com/dirc/expenses/internal/models"
	"gopkg.in/yaml.v3"
)

// LoadTransactionTypes loads yaml file with transactiontypes.
func LoadTransactionTypes(configPath string) ([]models.TransactionType, error) {
	// if !isSafePath(configPath) {
	// 	return fmt.Errorf("invalid file path: %s", configPath)
	// }
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config models.TransactionTypeConfig

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config.TransactionTypes, nil
}

// EnrichTransactions add transaction types for each transaction.
func EnrichTransactions(transactions []models.Transaction, types []models.TransactionType) []models.Transaction {
	for i, tx := range transactions {
		for _, tt := range types {
			if matchesSearchValues(tx, tt.SearchValues) {
				transactions[i].TransactionType = tt.Type

				break
			}
		}
	}

	return transactions
}

func matchesSearchValues(tx models.Transaction, searchValues map[string][]string) bool {
	for field, patterns := range searchValues {
		var value string

		switch field {
		case "naamTegenrekening":
			value = tx.NaamTegenrekening
		case "tegenrekening":
			value = tx.Tegenrekening
		case "omschrijving":
			value = tx.Omschrijving
		default:
			continue
		}

		for _, pattern := range patterns {
			// Handle wildcards
			if strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*") {
				// Match any substring
				if strings.Contains(strings.ToUpper(value), strings.ToUpper(strings.Trim(pattern, "*"))) {
					return true
				}
			} else if strings.HasPrefix(pattern, "*") {
				// Match suffix
				if strings.HasSuffix(strings.ToUpper(value), strings.ToUpper(strings.TrimPrefix(pattern, "*"))) {
					return true
				}
			} else if strings.HasSuffix(pattern, "*") {
				// Match prefix
				if strings.HasPrefix(strings.ToUpper(value), strings.ToUpper(strings.TrimSuffix(pattern, "*"))) {
					return true
				}
			} else {
				// Exact match
				if strings.EqualFold(value, pattern) {
					return true
				}
			}
		}
	}

	return false
}
