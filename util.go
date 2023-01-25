package main

import (
	"os"
	"fmt"

)

func cleanDB(dbName string) {
    err := os.Remove(dbName)
    if err != nil {
        if os.IsNotExist(err) {
            fmt.Println("File does not exist")
        } else {
            panic(err)
        }
    }
    fmt.Printf("File %s deleted successfully!\n", dbName)
}
