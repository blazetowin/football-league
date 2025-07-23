package main

import (
	"log"
	"net/http" 						

	"go-football-league/internal/api/routes"													
	"go-football-league/internal/repository"
)

func main() {
	// Initialize and establish the database connection
	repository.Connect()

	// Set up all API routes and return the configured router
	router := routes.SetupRouter()

	// Start the HTTP server on port 8080
	log.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
