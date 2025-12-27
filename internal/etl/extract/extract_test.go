package extract

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadCSV(t *testing.T) {
	// Create a temporary CSV file for testing
	csvContent := `"Boekdatum","Rekeningnummer","Bedrag","Debet / Credit","Naam tegenrekening","Tegenrekening","Code","Omschrijving","Saldo na boeking"
"07-12-2022","NL99TRIO0123456789","58,08","Debet","","","XS","CCV*Auto-Veer \DIRC","93,95"
"08-12-2022","NL99TRIO0123456789","12,80","Debet","","","XS","Bakkerij A","1.281,15"`

	tmpFile, err := os.CreateTemp(t.TempDir(), "test*.csv")
	require.NoError(t, err)

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(csvContent)
	require.NoError(t, err)
	err = tmpFile.Close()
	require.NoError(t, err)

	// Call the function
	transactions, err := ReadCSV(tmpFile.Name())
	require.NoError(t, err)

	// Assertions
	assert.Len(t, transactions, 2)

	// Check first transaction
	assert.Equal(t, "NL99TRIO0123456789", transactions[0].Rekeningnummer)
	assert.InEpsilon(t, 58.08, transactions[0].Bedrag, 0.001)
	assert.Equal(t, "Debet", transactions[0].DebetCredit)
	assert.Equal(t, "CCV*Auto-Veer \\DIRC", transactions[0].Omschrijving)
	assert.InEpsilon(t, 93.95, transactions[0].SaldoNaBoeking, 0.001)

	// Check second transaction
	assert.Equal(t, "NL99TRIO0123456789", transactions[1].Rekeningnummer)
	assert.InEpsilon(t, 12.80, transactions[1].Bedrag, 0.001)
	assert.Equal(t, "Bakkerij A", transactions[1].Omschrijving)
	assert.InEpsilon(t, 1281.15, transactions[1].SaldoNaBoeking, 0.001)

	// Check date parsing
	expectedDate, _ := time.Parse("02-01-2006", "07-12-2022")
	assert.Equal(t, expectedDate, transactions[0].Boekdatum)
}

func TestConvertAmountToFloat(t *testing.T) {
	// Define test cases
	tests := []struct {
		input    string  // Input string to convert
		expected float64 // Expected float64 value
		wantErr  bool    // Whether an error is expected
	}{
		{"1.518,01", 1518.01, false},        //nolint:golines // European format with dots as thousand separators and comma as decimal
		{"1518,01", 1518.01, false},         // European-style without thousand separator
		{"1.000,00", 1000.00, false},        // European format for 1000.00
		{"0,99", 0.99, false},               // European format small decimal
		{"invalid", 0, true},                // Invalid input
		{"", 0, true},                       // Empty string
		{"1.234.518,01", 1234518.01, false}, // Double decimal separator
	}

	// Iterate through test cases
	for _, tt := range tests {
		// Run the function
		result, err := convertAmountToFloat(tt.input)

		// Verify output
		if tt.wantErr {
			require.Error(t, err, "expected an error for input: '%s'", tt.input)
		} else {
			require.NoError(t, err, "did not expect an error for input: '%s'", tt.input)
			assert.InEpsilon(t, tt.expected, result, 0.0001, "unexpected result for input: '%s'", tt.input) // Tolerance for float comparison
		}
	}
}
