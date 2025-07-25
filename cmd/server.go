package main

import (
	"log"
	"net/http"

	"go-football-league/internal/api/routes"
	storage "go-football-league/internal/repository" // Database connection and schema loader
)

func main() {
	// Initialize the SQLite database and execute the schema setup
	storage.Connect()

	// Set up and return the router with all registered API endpoints
	router := routes.SetupRouter()

	// Start the HTTP server on port 8080
	log.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
