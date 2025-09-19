package models

import "time"

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

type TransactionTypeConfig struct {
	TransactionTypes []TransactionType `yaml:"TransactionTypes"`
}

type TransactionType struct {
	Type         string              `yaml:"type"`
	SearchValues map[string][]string `yaml:"SearchValues"`
}
