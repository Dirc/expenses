
# Expenses

Give clarity in your bank expenses.

## Run

Generate report from csv

```shell

./expenses \
  --csv=testdata/transactions.csv \
  --yaml=testdata/transactiontypes.yaml \
  --period=3y


```

## Development setup

```shell

go mod init
go mod tidy

go test ./... -v
golangci-lint run # Install golangci-lint from https://golangci-lint.run/usage/install

go build -ldflags="-X 'main.Version=$(git describe --tags --always)'" -o expenses cmd/main.go

```

## CSV Columns

```shell
Boekdatum	Rekeningnummer	Bedrag	Debet / Credit	Naam tegenrekening	Tegenrekening	Code	Omschrijving	Saldo na boeking
```

## ToDo

### V2.0

- [x] ETL project structure
- [x] Extract, Transform + unit tests
- [x] Reports: `3m` for 3 months, `2y` for 2 years

### v2.1

- [x] report: all untyped transactions
- [x] CLI
- [x] linting

### v2.2

- [ ] Seperate report and untyped (+ cli commands)
- [ ] show transactions per type
- [ ] rename: type -> category?

### future ideas

- [ ] Load to sqlite (or duckdb?) -> nice, but so far the data is small enough to generate on the fly
- [ ] incremental updates -> idem
- [ ] UI
- [ ] Variable csv columns

### MVP (old)

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

#### nice to haves

- [ ] provide api with Gin
- [ ] support duckdb as backend

