package league

import (
	"fmt"

	storage "go-football-league/internal/repository"
)

// PlayWeek runs the simulation process for a specific week.
// It first checks whether the week has already been simulated to prevent duplicate execution.
// If not played it creates fixtures and simulates the match results.
// Returns an error if any step fails.
func PlayWeek(week int) error {
	if played, err := weekAlreadyPlayed(week); err != nil {
		return fmt.Errorf("Failed to check if week was already played: %v", err)
	} else if played {
		fmt.Printf("Week %d already played. Skipping.\n", week)
		return nil
	}

	fmt.Printf("Generating fixtures for week %d...\n", week)

	if err := GenerateWeeklyMatches(week); err != nil {
		return fmt.Errorf("Failed to generate weekly matches: %v", err)
	}

	fmt.Printf("Simulating results for week %d...\n", week)

	if err := SimulateScores(week); err != nil {
		return fmt.Errorf("Failed to simulate match scores: %v", err)
	}

	fmt.Printf("Week %d simulation completed.\n", week)
	return nil
}

// weekAlreadyPlayed determines whether the given week already has recorded results.
// It queries the database for matches with non-null score values.
// Returns true if the week has already been played.
func weekAlreadyPlayed(week int) (bool, error) {
	var count int
	err := storage.DB.QueryRow(`
		SELECT COUNT(*) FROM matches 
		WHERE week = ? AND home_goals IS NOT NULL AND away_goals IS NOT NULL
	`, week).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// PrintMatchesOfWeek prints the match results or fixtures for the given week.
// If match scores are present, it displays them; otherwise, it shows placeholders.
func PrintMatchesOfWeek(week int) error {
	rows, err := storage.DB.Query(`
		SELECT m.id, t1.name, t2.name, m.home_goals, m.away_goals
		FROM matches m
		JOIN teams t1 ON m.home_team_id = t1.id
		JOIN teams t2 ON m.away_team_id = t2.id
		WHERE m.week = ?
	`, week)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Printf("Match results for week %d:\n", week)
	for rows.Next() {
		var matchID int
		var homeTeam, awayTeam string
		var homeGoals, awayGoals *int
		if err := rows.Scan(&matchID, &homeTeam, &awayTeam, &homeGoals, &awayGoals); err != nil {
			return err
		}

		score := "vs" // Default text if match has not been played
		if homeGoals != nil && awayGoals != nil {
			score = fmt.Sprintf("%d-%d", *homeGoals, *awayGoals)
		}

		fmt.Printf("  Match %d: %s %s %s\n", matchID, homeTeam, score, awayTeam)
	}
	return nil
}
