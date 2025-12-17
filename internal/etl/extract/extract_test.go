package extract

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadCSV(t *testing.T) {
	// Create a temporary CSV file for testing
	csvContent := `"Boekdatum","Rekeningnummer","Bedrag","Debet / Credit","Naam tegenrekening","Tegenrekening","Code","Omschrijving","Saldo na boeking"
"07-12-2022","NL99TRIO0123456789","58,08","Debet","","","XS","CCV*Auto-Veer \DIRC","93,95"
"08-12-2022","NL99TRIO0123456789","12,80","Debet","","","XS","Bakkerij A","81,15"`

	tmpFile, err := os.CreateTemp(t.TempDir(), "test*.csv")
	assert.NoError(t, err)
	// require.noError(t, err)

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(csvContent)
	assert.NoError(t, err)
	err = tmpFile.Close()
	assert.NoError(t, err)

	// Call the function
	transactions, err := ReadCSV(tmpFile.Name())
	assert.NoError(t, err)

	// Assertions
	assert.Len(t, transactions, 2)

	// Check first transaction
	assert.Equal(t, "NL99TRIO0123456789", transactions[0].Rekeningnummer)
	assert.Equal(t, "Debet", transactions[0].DebetCredit)
	assert.Equal(t, "CCV*Auto-Veer \\DIRC", transactions[0].Omschrijving)
	assert.InEpsilon(t, 58.08, transactions[0].Bedrag, 0.001)
	assert.InEpsilon(t, 93.95, transactions[0].SaldoNaBoeking, 0.001)

	// Check second transaction
	assert.Equal(t, "NL99TRIO0123456789", transactions[1].Rekeningnummer)
	assert.Equal(t, "Bakkerij A", transactions[1].Omschrijving)
	assert.InEpsilon(t, 12.80, transactions[1].Bedrag, 0.001)
	assert.InEpsilon(t, 81.15, transactions[1].SaldoNaBoeking, 0.001)

	// Check date parsing
	expectedDate, _ := time.Parse("02-01-2006", "07-12-2022")
	assert.Equal(t, expectedDate, transactions[0].Boekdatum)
}
