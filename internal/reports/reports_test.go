package reports

import (
	"testing"
	"time"

	"github.com/dirc/expenses/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGenerateReport(t *testing.T) {
	// Mock transactions
	now := time.Now()
	transactions := []models.Transaction{
		{
			Boekdatum:       now.AddDate(0, -1, 0), // 1 month ago
			Bedrag:          100.00,
			TransactionType: "Boodschappen",
		},
		{
			Boekdatum:       now.AddDate(0, -2, 0), // 2 months ago
			Bedrag:          200.00,
			TransactionType: "Auto",
		},
		{
			Boekdatum:       now.AddDate(0, -4, 0), // 4 months ago (should be excluded for 3m)
			Bedrag:          300.00,
			TransactionType: "Boodschappen",
		},
		{
			Boekdatum:       now.AddDate(0, 0, 0), // Today
			Bedrag:          400.00,
			TransactionType: "",
		},
	}

	t.Run("Monthly Report", func(t *testing.T) {
		report, err := GeneratePeriodicReport(transactions, "3m")
		assert.NoError(t, err)
		assert.Len(t, report, 3) // Only 3 periods within last 3 months

		// Check amounts
		for period, types := range report {
			if period == now.AddDate(0, -1, 0).Format("2006-01") {
				assert.InEpsilon(t, 100.00, types["Boodschappen"], 0.001)
			}

			if period == now.AddDate(0, 0, 0).Format("2006-01") {
				assert.InEpsilon(t, 400.00, types["Untyped"], 0.001)
			}
		}
	})

	t.Run("Yearly Report", func(t *testing.T) {
		report, err := GeneratePeriodicReport(transactions, "1y")
		assert.NoError(t, err)
		assert.Len(t, report, 1) // All transactions within last year

		// Check amounts
		for period, types := range report {
			if period == now.Format("2006") {
				assert.InEpsilon(t, 400.00, types["Boodschappen"], 0.001)
				assert.InEpsilon(t, 200.00, types["Auto"], 0.001)
				assert.InEpsilon(t, 400.00, types["Untyped"], 0.001)
			}
		}
	})

	t.Run("Invalid Duration", func(t *testing.T) {
		_, err := GeneratePeriodicReport(transactions, "3x")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid duration format")
	})

	t.Run("Empty Transactions", func(t *testing.T) {
		report, err := GeneratePeriodicReport([]models.Transaction{}, "3m")
		assert.NoError(t, err)
		assert.Empty(t, report)
	})
}
