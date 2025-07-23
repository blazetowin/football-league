package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // SQLite driver import 
)

var DB *sql.DB // Global database connection handle

// Connect initializes the SQLite database connection and executes schema setup.
// It reads the schema SQL file and runs its statements to create required tables.
// If any step fails, the application logs the error and terminates.
func Connect() {
	var err error

	// Open or create the SQLite database file
	DB, err = sql.Open("sqlite3", "./league.db")
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Load SQL schema from file
	schema, err := os.ReadFile("internal/migration/schema.sql")
	if err != nil {
		log.Fatal("Could not read schema.sql file:", err)
	}

	// Execute the schema SQL to set up tables
	_, err = DB.Exec(string(schema))
	if err != nil {
		log.Fatal("Failed to execute schema:", err)
	}

	fmt.Println("Database connection established and schema applied successfully.")
}
