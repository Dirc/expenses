
package main

import (
	"database/sql"
	"log"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func createTable(db *sql.DB) error {
    _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY,
        name TEXT NOT NULL,
        age INTEGER NOT NULL,
        sexe TEXT
    )`)
    if err != nil {
        return fmt.Errorf("Error creating table: %v", err)
    }
    return nil
}

func insertData(name string, age int, db *sql.DB) error {
    _, err := db.Exec(`INSERT INTO users (name, age) VALUES (?, ?)`, name, age)
    if err != nil {
        return fmt.Errorf("Error inserting data: %v", err)
    }
    return nil
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

func generateSexe(transTyp TransactionType) {
    db, err := sql.Open("sqlite3", "./dev.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    for column := range transTyp.SearchValues {
        fmt.Println(column)
        for _, value := range transTyp.SearchValues[column] {
            // Build query
            query := "UPDATE users SET sexe = '" + transTyp.Type + "' WHERE " + column + " LIKE '%" + value + "%'"
            fmt.Println(query)

            // Prepare the statement
            stmt, err := db.Prepare(query)
            if err != nil {
                panic(err)
            }
            defer stmt.Close()

            // Execute the query
            result, err := stmt.Exec()
            if err != nil {
                panic(err)
            }

            // Get the number of affected rows
            affected, err := result.RowsAffected()
            if err != nil {
                panic(err)
            }
            fmt.Printf("Number of affected rows: %d\n", affected)
        }
    }

    fmt.Println("Values inserted successfully!")
}

func returnWomen() {
    db, err := sql.Open("sqlite3", "./dev.db")
    rows, err := db.Query("SELECT name, age, sexe FROM users WHERE sexe = 'women'")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var name string
        var age int
        var sexe string
        if err := rows.Scan(&name, &age, &sexe); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%s is a %d years old %s\n", name, age, sexe)
    }

    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
}

func printTable() {
    db, err := sql.Open("sqlite3", "./dev.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    rows, err := db.Query("SELECT id, name, age, sexe FROM users")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var id int
        var name string
        var age int
        var sexe sql.NullString
        if err := rows.Scan(&id, &name, &age, &sexe); err != nil {
            log.Fatal(err)
        }
        if sexe.Valid {
            fmt.Printf("id: %d, name: %s, age: %d, sexe: %s\n", id, name, age, sexe.String)
        } else {
            fmt.Printf("id: %d, name: %s, age: %d, sexe: unknown\n", id, name, age)
        }
    }

    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
}

func loopOverMap() {

    m := map[string][]string{
        "key0": []string{"category"},
        "key1": []string{"Alice", "Bob", "Charlie"},
        "key2": []string{"Dave", "Eve", "Frank"},
        "key3": []string{"Gary", "Hannah", "Ian"},
    }

    category := m["key0"][0]
    delete(m, "key0")

    // Loop over key-value pairs
    for k, _ := range m {
        for _, v := range m[k] {
            fmt.Printf("Key: %s, Value: %s, Category: %s\n", k, v, category)
        }
    }
}

func foo() int {
    return 42
}

type TransactionType struct {
    Type            string
    SearchValues    map[string][]string
}

func main() {
  db, err := sql.Open("sqlite3", "./dev.db")
  if err != nil {
      panic(err)
  }
  defer db.Close()

  err = createTable(db)
  if err != nil {
      panic(err)
  }

  insertData("Mijke", 37, db)
  insertData("Resa", 5, db)
  insertData("Maud", 7, db)
  insertData("Eric", 37, db)

  returnAge(18)

  transTypWomen := TransactionType{
    Type: "women",
    SearchValues: map[string][]string{
        "name": []string{"Mijke", "Maud", "Resa"},
    },
  }

  fmt.Println(transTypWomen)
  fmt.Println(transTypWomen.Type)
  fmt.Println(transTypWomen.SearchValues)
  fmt.Println(transTypWomen.SearchValues["name"][0])
  // Print keys of map
  for key := range transTypWomen.SearchValues {
    fmt.Println(key)
  }


  generateSexe(transTypWomen)
  returnWomen()

  printTable()

  cleanDB()



//   loopOverMap()

}