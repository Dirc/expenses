
package main

import (
	"database/sql"
	"log"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func createTable() {
    db, err := sql.Open("sqlite3", "./dev.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY,
        name TEXT NOT NULL,
        age INTEGER NOT NULL
    )`)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Table created successfully!")
}

func insertData(name string, age int) {
    db, err := sql.Open("sqlite3", "./dev.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    _, err = db.Exec(`INSERT INTO users (name, age) VALUES (?, ?)`, name, age)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Row inserted successfully!")
}

func returnAge(age int) {
    db, err := sql.Open("sqlite3", "./dev.db")
    rows, err := db.Query("SELECT name, age FROM users WHERE age > $1", age)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var name string
        var age int
        if err := rows.Scan(&name, &age); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%s is %d years old\n", name, age)
    }

    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
}

func cleanDB() {
    err := os.Remove("./dev.db")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("File deleted successfully!")
}

func main() {
  createTable()
  insertData("Mijke", 37)
  insertData("Resa", 5)
  insertData("Maud", 7)
  insertData("Eric", 37)
  returnAge(18)
  cleanDB()
}