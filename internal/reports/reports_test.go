package reports

import (
	"testing"
	"time"

	"github.com/dirc/expenses/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			Boekdatum:       now.AddDate(0, -1, 0), // 1 months ago
			Bedrag:          200.00,
			TransactionType: "Auto",
		},
		{
			Boekdatum:       now.AddDate(0, -2, 0), // 2 months ago (should be excluded for 1m)
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
		report, err := GeneratePeriodicReport(transactions, "1m")
		require.NoError(t, err)
		assert.Len(t, report, 3) // Only 3 periods within last 3 months. Todo: mock time, this fails in january since it uses real time

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
		require.NoError(t, err)
		assert.Len(t, report, 1) // All transactions within last year. Todo: mock time, this fails in january, february since it uses real time

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
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid duration format")
	})

	t.Run("Empty Transactions", func(t *testing.T) {
		report, err := GeneratePeriodicReport([]models.Transaction{}, "3m")
		require.NoError(t, err)
		assert.Empty(t, report)
	})
}

func TestFilterByPeriod(t *testing.T) {
	now := time.Now()
	transactions := []models.Transaction{
		{
			Boekdatum:       now.AddDate(0, -1, 0),
			Omschrijving:    "Recent Untyped",
			TransactionType: "", // Untyped
		},
		{
			Boekdatum:       now.AddDate(0, -6, 0),
			Omschrijving:    "Old Untyped",
			TransactionType: "", // Untyped
		},
	}

	t.Run("Filters untyped transactions by date", func(t *testing.T) {
		filtered, err := FilterByPeriod(transactions, "3m")
		require.NoError(t, err)

		// Only the recent untyped transaction should remain
		assert.Len(t, filtered, 1)
		assert.Equal(t, "Recent Untyped", filtered[0].Omschrijving)
	})

	t.Run("Invalid duration format", func(t *testing.T) {
		_, err := FilterByPeriod(transactions, "10d") // 'd' is not supported
		assert.Error(t, err)
	})
}
