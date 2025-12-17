// Package main creates the CLI.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/dirc/expenses/internal/etl/extract"
	"github.com/dirc/expenses/internal/etl/transform"
	"github.com/dirc/expenses/internal/reports"
)

func main() {
	// Define CLI flags
	csvPath := flag.String("csv", "", "Path to the CSV file containing transactions")
	yamlPath := flag.String("yaml", "", "Path to the YAML file containing transaction types")
	reportPeriod := flag.String("period", "", "Report period (e.g., '3m' for 3 months, '4y' for 4 years)")

	// Parse flags
	flag.Parse()

	// Validate required flags
	if *csvPath == "" || *yamlPath == "" || *reportPeriod == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Extract transactions from CSV
	transactions, err := extract.ReadCSV(*csvPath)
	if err != nil {
		log.Fatalf("Failed to read CSV: %v", err)
	}

	// Load transaction types from YAML
	types, err := transform.LoadTransactionTypes(*yamlPath)
	if err != nil {
		log.Fatalf("Failed to load transaction types: %v", err)
	}

	// Enrich transactions with types
	transactions = transform.EnrichTransactions(transactions, types)

	// Generate and print reports
	report, err := reports.GeneratePeriodicReport(transactions, *reportPeriod)
	if err != nil {
		log.Fatalf("Failed to generate report: %v", err)
	}

	// Print the report
	reports.PrintPeriodicReport(report, log.Printf("Report for %s", reportPeriod))
	reports.GenerateUntypedReport(transactions)
}
