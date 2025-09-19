package main

import (
	"log"

	"github.com/dirc/expenses/internal/etl/extract"
	"github.com/dirc/expenses/internal/etl/transform"
	"github.com/dirc/expenses/internal/reports"
)

func main() {
	// Path to your CSV file
	csvPath := "testdata/transactions.csv"
	// Path to your YAML config
	configPath := "testdata/transactiontypes.yaml"

	// Extract transactions from CSV
	transactions, err := extract.ReadCSV(csvPath)
	if err != nil {
		log.Fatalf("Failed to read CSV: %v", err)
	}

	// Load transaction types from YAML
	types, err := transform.LoadTransactionTypes(configPath)
	if err != nil {
		log.Fatalf("Failed to load transaction types: %v", err)
	}

	// Enrich transactions with types
	transactions = transform.EnrichTransactions(transactions, types)

	// Generate and print reports
	monthReport, err := reports.GenerateReport(transactions, "3m") // Last 3 months
	if err != nil {
		log.Fatalf("Failed to generate monthly report: %v", err)
	}
	yearReport, err := reports.GenerateReport(transactions, "4y") // Last 4 years
	if err != nil {
		log.Fatalf("Failed to generate yearly report: %v", err)
	}

	reports.PrintReport(monthReport, "Monthly Report (Last 3 Months)")
	reports.PrintReport(yearReport, "Yearly Report (Last 4 Years)")
}
