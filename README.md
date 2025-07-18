
# Expenses

Give clarity in your bank expenses.

## Run

Import csv into sqlite DB.

```shell
go mod init
go mod tidy

go run *.go

go test *.go

```

## Columns

```shell
Boekdatum	Rekeningnummer	Bedrag	Debet / Credit	Naam tegenrekening	Tegenrekening	Code	Omschrijving	Saldo na boeking
```

## ToDo

### MVP

- [x] add column "transactionType" e.g. boodschappen, vakantie, vasteLasten, fun
- [x] create db connection in main and use db as input for other funcs. (as in dev.go createTable)
- [x] make generateTransactionType more general
  - [x] general for loop sql
  - [x] struct as input
- [ ] unit tests
  - [x] create-db
  - [ ] transactionTypes
- [x] print total amount per transactionType
- [x] print pretty with padding
- [x] define multiple transaction structs
- [ ] run as cli
  - [ ] init
  - [ ] load transactionTypes (from yaml?) 

- [ ] transactionTypesTotal:
  - [ ] add time: "from" till "when" the data is captured
  - [ ] add column for expenses per month
- [x] printUnknown: print table for all unknown

- [ ] store struct in seperate table?
  - [ ] method to add/rm search items

- [ ] make importing csv more general
  - input: csv, list of columns?

### nice to haves

- [ ] provide api with Gin
- [ ] support duckdb as backend

