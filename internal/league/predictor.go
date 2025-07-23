package league

import (
	"fmt"
	"math"
	"sort"
	"go-football-league/internal/domain"
)
// PrintChampionshipPredictions displays the league title chances for each team based on the current standings of the given week.
// It only runs if the league has progressed to at least week 4.
func PrintChampionshipPredictions(week int, table []models.LeagueTableRow) {
	if week < 4 {
		// Not enough data to calculate predictions
		return
	}

	fmt.Printf("\nChampionship Predictions - Week %d:\n", week)

	// Calculate total points in the table
	totalPoints := 0
	for _, t := range table {
		totalPoints += t.Points
	}

	// If total points is 0, assign equal chances to all teams
	if totalPoints == 0 {
		for _, t := range table {
			fmt.Printf("%-15s 25%%\n", t.TeamName)
		}
		return
	}

	// Compute each team's chance based on their points
	type Prediction struct {
		Team string
		Rate float64
	}

	var preds []Prediction
	for _, t := range table {
		p := (float64(t.Points) / float64(totalPoints)) * 100
		preds = append(preds, Prediction{
			Team: t.TeamName,
			Rate: p,
		})
	}

	// Sort predictions in descending order
	sort.Slice(preds, func(i, j int) bool {
		return preds[i].Rate > preds[j].Rate
	})

	// Display rounded predictions 
	for _, p := range preds {
		fmt.Printf("%-15s %.0f%%\n", p.Team, math.Round(p.Rate))
	}
}
