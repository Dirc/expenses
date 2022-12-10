
# Expenses

Order your bank expenses

## Run

Import csv into sqlite DB.

```shell
go mod init
go mod tidy

go run main.go

echo "select * from expenses;" | sqlite3 expenses.db

```

## SQL

Print expenses table

```shell
echo "select * from expenses;" | sqlite3 expenses.db
```

## Columns

```shell
Boekdatum	Rekeningnummer	Bedrag	Debet / Credit	Naam tegenrekening	Tegenrekening	Code	Omschrijving	Saldo na boeking
```
