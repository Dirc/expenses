
chatgpt

```golang
package main

import (
"database/sql"
"fmt"

    _ "github.com/go-sql-driver/mysql"
)

func main() {
// Open a connection to the database
db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/database")
if err != nil {
panic(err.Error())
}
defer db.Close()

    // Insert a new record into the "users" table
    result, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", "John", "john@example.com")
    if err != nil {
        panic(err.Error())
    }

    // Get the ID of the inserted record
    id, err := result.LastInsertId()
    if err != nil {
        panic(err.Error())
    }
    fmt.Println("Inserted record ID:", id)
}
```