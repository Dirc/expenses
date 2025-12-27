package transform

import (
	"os"
	"testing"

	"github.com/dirc/expenses/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchesSearchValues(t *testing.T) {
	tx := models.Transaction{
		NaamTegenrekening: "PANIC",
		Tegenrekening:     "NL123",
		Omschrijving:      "Bakkerij De Vries",
	}

	// Test exact match
	assert.True(t, matchesSearchValues(tx, map[string][]string{
		"naamTegenrekening": {"PANIC"},
	}))
	assert.False(t, matchesSearchValues(tx, map[string][]string{
		"naamTegenrekening": {"RABO"},
	}))

	// Test wildcard (*) matching
	assert.True(t, matchesSearchValues(tx, map[string][]string{
		"omschrijving": {"*Bakkerij*"},
	}))
	assert.True(t, matchesSearchValues(tx, map[string][]string{
		"omschrijving": {"*De Vries"},
	}))
	assert.True(t, matchesSearchValues(tx, map[string][]string{
		"omschrijving": {"Bakkerij*"},
	}))
	assert.False(t, matchesSearchValues(tx, map[string][]string{
		"omschrijving": {"*Supermarkt*"}, //nolint:misspell
	}))

	// Test case insensitivity
	assert.True(t, matchesSearchValues(tx, map[string][]string{
		"naamTegenrekening": {"panic"},
	}))
	assert.True(t, matchesSearchValues(tx, map[string][]string{
		"omschrijving": {"*bakkerij*"},
	}))
}

func TestEnrichTransactions(t *testing.T) {
	transactions := []models.Transaction{
		{
			NaamTegenrekening: "PANIC",
			Tegenrekening:     "NL123",
			Omschrijving:      "Bakkerij De Vries",
		},
		{
			NaamTegenrekening: "RABO",
			Tegenrekening:     "NL456",
			Omschrijving:      "Auto-Veer Betaalautomaat",
		},
	}

	types := []models.TransactionType{
		{
			Type: "Boodschappen",
			SearchValues: map[string][]string{
				"naamTegenrekening": {"PANIC"},
				"omschrijving":      {"*Bakkerij*"},
			},
		},
		{
			Type: "Auto",
			SearchValues: map[string][]string{
				"omschrijving": {"*Auto-Veer*"},
			},
		},
	}

	enriched := EnrichTransactions(transactions, types)

	assert.Equal(t, "Boodschappen", enriched[0].TransactionType)
	assert.Equal(t, "Auto", enriched[1].TransactionType)
}

func TestLoadTransactionTypes(t *testing.T) {
	// Create a temporary YAML file for testing
	yamlContent := `
transactionTypes:
- type: Boodschappen
  searchValues:
    naamTegenrekening: ["PANIC"]
    omschrijving: ["*Bakkerij*"]

- type: Auto
  searchValues:
    omschrijving: ["*Auto-Veer*"]
    `

	tmpFile, err := os.CreateTemp(t.TempDir(), "test*.yaml")
	require.NoError(t, err)

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(yamlContent)
	require.NoError(t, err)
	err = tmpFile.Close()
	require.NoError(t, err)

	types, err := LoadTransactionTypes(tmpFile.Name())
	require.NoError(t, err)
	assert.Len(t, types, 2)
	assert.Equal(t, "Boodschappen", types[0].Type)
	assert.Equal(t, "Auto", types[1].Type)
}
