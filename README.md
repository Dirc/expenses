
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

# All debet transactions
echo "select * from expenses WHERE debetCredit = 'Debet';" | sqlite3 expenses.db

# All transaction on 8th of the month
echo "select * from expenses WHERE boekdatum LIKE '08-%';" | sqlite3 expenses.db

# All Debet transaction on 8th of the month
echo "select * from expenses WHERE boekdatum LIKE '08-%' AND debetCredit = 'Debet';" | sqlite3 expenses.db

# Column names
echo "PRAGMA table_info(expenses)" | sqlite3 expenses.db

# Sum of transactions
echo "SELECT SUM(bedrag) FROM expenses" | sqlite3 expenses.db

# Sum of debet transactions
echo "SELECT SUM(bedrag) FROM expenses WHERE debetCredit = 'Debet'" | sqlite3 expenses.db

# Update transactionType based on naamTegenrekening
echo "UPDATE expenses SET transactionType = 'booschappen' WHERE naamTegenrekening = 'PICNIC BY BUCKAROO';" | sqlite3 expenses.db 

```

## Columns

```shell
Boekdatum	Rekeningnummer	Bedrag	Debet / Credit	Naam tegenrekening	Tegenrekening	Code	Omschrijving	Saldo na boeking
```

## ToDo

- [x] add column "transactionType" e.g. boodschappen, vakantie, vasteLasten, fun
- [ ] make generateTransactionType more general (map as input)
  - [x] general for loop sql
  - [ ] map as input
  - [ ] define multiple maps
    - store map in seperate table?
- [ ] print total amount per transactionType
- [ ] unit tests