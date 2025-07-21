package storage

import (
	"database/sql" // Importing the database/sql package for database operations
	"fmt" // Importing fmt for formatted I/O
	"log" // Importing log for logging errors
	"os" // Importing os for file operations

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var DB *sql.DB

// Connect establishes a connection to the database and sets up the schema
func Connect() {
	var err error

	// Open SQLite database (creates if not exists)
	DB, err = sql.Open("sqlite3", "./league.db")
	if err != nil {
		log.Fatal("❌ Failed to connect to the database:", err)
	}

	// Read the schema.sql file
	schema, err := os.ReadFile("sql/schema.sql")
	if err != nil {
		log.Fatal("❌ Could not read schema.sql:", err)
	}

	// Execute the SQL commands
	_, err = DB.Exec(string(schema))
	if err != nil {
		log.Fatal("❌ Failed to execute schema:", err)
	}

	fmt.Println("✅ Database connection and table setup successful.")
}