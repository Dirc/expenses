package reports

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dirc/expenses/internal/models"
)

// GenerateReport generates a report for the last N periods (months/years).
// durationStr: e.g., "3m" for 3 months, "4y" for 4 years
func GenerateReport(
	transactions []models.Transaction,
	durationStr string,
) (map[string]map[string]float64, error) {
	// Parse duration string (e.g., "3m" or "4y")
	duration, err := strconv.Atoi(durationStr[:len(durationStr)-1])
	if err != nil {
		return nil, fmt.Errorf("invalid duration: %v", err)
	}

	isMonth := strings.HasSuffix(durationStr, "m")
	isYear := strings.HasSuffix(durationStr, "y")
	if !isMonth && !isYear {
		return nil, fmt.Errorf("invalid duration format, use e.g. '3m' or '4y'")
	}

	now := time.Now()
	report := make(map[string]map[string]float64)

	// Calculate the cutoff date
	var cutoff time.Time
	if isMonth {
		cutoff = now.AddDate(0, -duration, 0)
	} else {
		cutoff = now.AddDate(-duration, 0, 0)
	}

	for _, tx := range transactions {
		// Skip transactions older than the cutoff
		if tx.Boekdatum.Before(cutoff) {
			continue
		}

		// Format the period key
		var periodKey string
		if isMonth {
			periodKey = tx.Boekdatum.Format("2006-01")
		} else {
			periodKey = tx.Boekdatum.Format("2006")
		}
		txType := tx.TransactionType
		if txType == "" {
			txType = "Untyped"
		}

		if _, ok := report[periodKey]; !ok {
			report[periodKey] = make(map[string]float64)
		}
		report[periodKey][txType] += tx.Bedrag
	}

	return report, nil
}

// PrintReport prints the report in a readable format.
func PrintReport(report map[string]map[string]float64, title string) {
	fmt.Printf("\n=== %s ===\n", title)
	for period, types := range report {
		fmt.Printf("\nPeriod: %s\n", period)
		for txType, amount := range types {
			fmt.Printf("  %s: €%.2f\n", txType, amount)
		}
	}
}
