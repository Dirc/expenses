// Package models for all models.
package models

import "time"

// Transaction of a bank account with added column: TransactionType.
type Transaction struct {
	Boekdatum         time.Time `csv:"Boekdatum"`
	Rekeningnummer    string    `csv:"Rekeningnummer"`
	Bedrag            float64   `csv:"Bedrag"`
	DebetCredit       string    `csv:"Debet / Credit"`
	NaamTegenrekening string    `csv:"Naam tegenrekening"`
	Tegenrekening     string    `csv:"Tegenrekening"`
	Code              string    `csv:"Code"`
	Omschrijving      string    `csv:"Omschrijving"`
	SaldoNaBoeking    float64   `csv:"Saldo na boeking"`
	TransactionType   string    // Will be set during transform
}

// TransactionTypeConfig is a array of transactionTypes.
type TransactionTypeConfig struct {
	TransactionTypes []TransactionType `yaml:"transactionTypes"`
}

// TransactionType has a name and one or more search values so it can be matched on a transaction column.
type TransactionType struct {
	Type         string              `yaml:"type"`
	SearchValues map[string][]string `yaml:"searchValues"`
}
