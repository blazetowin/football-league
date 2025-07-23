package league

import (
	"fmt"
	models "go-football-league/internal/domain"
)

// PrintLeagueTableRows renders a formatted league table to the console.
// It displays each team's performance statistics in a tabular view.
func PrintLeagueTableRows(table []models.LeagueTableRow) {
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("%-15s %2s %2s %2s %2s %4s %4s %4s %4s\n", 
		"Team", "MP", "W", "D", "L", "GF", "GA", "GD", "Pts")
	fmt.Println("------------------------------------------------------------")

	// Print each row of the table with aligned columns
	for _, row := range table {
		fmt.Printf("%-15s %2d %2d %2d %2d %4d %4d %4d %4d\n",
			row.TeamName, row.Played, row.Wins, row.Draws, row.Losses,
			row.GoalsFor, row.GoalsAgainst, row.GoalDiff, row.Points)
	}

	fmt.Println("------------------------------------------------------------")
}
