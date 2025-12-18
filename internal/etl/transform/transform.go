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
	data, err := os.ReadFile(configPath) // #nosec G304
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
func EnrichTransactions(
	transactions []models.Transaction,
	types []models.TransactionType,
) []models.Transaction {
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
			upperValue := strings.ToUpper(value)

			// Handle wildcards
			switch {
			case strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*"):
				// Match any substring
				upperPattern := strings.ToUpper(strings.Trim(pattern, "*"))
				if strings.Contains(upperValue, upperPattern) {
					return true
				}
			case strings.HasPrefix(pattern, "*"):
				// Match suffix
				if strings.HasSuffix(upperValue, strings.ToUpper(strings.TrimPrefix(pattern, "*"))) {
					return true
				}
			case strings.HasSuffix(pattern, "*"):
				// Match prefix
				if strings.HasPrefix(upperValue, strings.ToUpper(strings.TrimSuffix(pattern, "*"))) {
					return true
				}
			default:
				// Exact match
				if strings.EqualFold(value, pattern) {
					return true
				}
			}
		}
	}

	return false
}
