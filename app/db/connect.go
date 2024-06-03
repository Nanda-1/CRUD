package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Connect function establishes a connection to the PostgreSQL database and returns the *sql.DB object
func Connect() *sql.DB {
	// Update the connection string with your actual database credentials
	connStr := "user=postgres password=admin dbname=blog sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Verify the connection with a Ping
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Successfully connected to the database")
	return db
}
