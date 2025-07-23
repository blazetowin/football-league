package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"go-football-league/internal/league"
	"go-football-league/internal/repository"
)

func main() {
	// Initialize the database and apply schema
	storage.Connect()
	fmt.Println("Database connection and schema setup complete.")

	// Create the fixture if it doesn't already exist
	err := league.CreateFixture()
	if err != nil {
		log.Fatalf("Failed to create fixture: %v", err)
	}

	// Simulate each of the 6 weeks in the league
	for week := 1; week <= 6; week++ {
		fmt.Printf("===== WEEK %d =====\n\n", week)

		// Generate matches for this week (if not already created)
		err := league.GenerateWeeklyMatches(week)
		if err != nil {
			log.Fatalf("Failed to generate matches for week %d: %v", week, err)
		}

		// Simulate scores for this week's matches
		err = league.SimulateScores(week)
		if err != nil {
			log.Fatalf("Failed to simulate scores for week %d: %v", week, err)
		}

		// Display match results
		fmt.Printf("\nMatch Results (Week %d):\n", week)
		league.PrintMatchesOfWeek(week)

		// Generate and display the updated league table
		fmt.Printf("\nLeague Standings (After Week %d):\n", week)
		table, err := league.GenerateLeagueTable(week)
		if err != nil {
			log.Fatalf("Failed to generate league table: %v", err)
		}
		league.PrintLeagueTableRows(table)

		// If week >= 4, print championship predictions
		// This will only run if the league has progressed to at least week 4
		if week >= 4 {
			league.PrintChampionshipPredictions(week, table)
		}

		// Pause between weeks for user input
		// This is optional but can help in observing the simulation step-by-step
		if week < 6 {
			fmt.Print("\nPress Enter to continue to the next week...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			fmt.Println()
		}
	}
}

