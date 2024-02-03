package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Replace the following connection string with your PostgreSQL database details
	connectionString := "user=root password=secret dbname=tutor_db sslmode=disable"

	// Open a connection to the PostgreSQL server
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the PostgreSQL server!")

	// Perform database operations here...
}
