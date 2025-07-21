package league
	


import (
	"fmt"
	"go-football-league/internal/models"
)

func PrintLeagueTableRows(table []models.LeagueTableRow) {
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("%-15s %2s %2s %2s %2s %4s %4s %4s %4s\n", "Takım", "O", "G", "B", "M", "AG", "YG", "+/-", "P")
	fmt.Println("------------------------------------------------------------")
	for _, row := range table {
		fmt.Printf("%-15s %2d %2d %2d %2d %4d %4d %4d %4d\n",
			row.TeamName, row.Played, row.Wins, row.Draws, row.Losses,
			row.GoalsFor, row.GoalsAgainst, row.GoalDiff, row.Points)
	}
	fmt.Println("------------------------------------------------------------")
}