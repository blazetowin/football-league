package league

import (
	"fmt" 	// Importing fmt for formatted I/O

	"go-football-league/internal/storage" // Importing storage for database operations
)

// PlayWeek plays one simulation week (match creation + result generation)
// This function is responsible for simulating a week in the football league.
// It first generates matches for the specified week and then simulates the scores for those matches.
// It returns an error if any step fails, allowing the caller to handle it appropriately.
func PlayWeek(week int) error {
	// Check if the week has already been played
	// Prevent re-playing the same week
	if played, err := weekAlreadyPlayed(week); err != nil {
		return fmt.Errorf("❌ Failed to check if week played: %v", err)
	} else if played {
		fmt.Printf("❌ Week %d already played. Skipping.\n", week)
		return nil
	}

	fmt.Printf("ℹ️ Generating matches for week %d...\n", week)
	// 4 takımdan 2 maç oluşturuyoruz
	err := GenerateWeeklyMatches(week)
	if err != nil {
		return fmt.Errorf("❌ Failed to generate matches: %v", err)
	}

	fmt.Println("📅 Week", week, "is being simulated...")
	err = SimulateScores(week)
	if err != nil {
		return fmt.Errorf("❌ Failed to simulate scores: %v", err)
	}

	fmt.Println("✅ Week", week, "completed.")
	return nil
}

// weekAlreadyPlayed checks if the week has already been played by verifying if any matches with scores exist.
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

	fmt.Printf("🏟️ Matches for week %d:\n", week)
	for rows.Next() {
		var matchID int
		var homeTeam, awayTeam string
		var homeGoals, awayGoals *int
		if err := rows.Scan(&matchID, &homeTeam, &awayTeam, &homeGoals, &awayGoals); err != nil {
			return err
		}
		score := "vs"
		if homeGoals != nil && awayGoals != nil {
			score = fmt.Sprintf("%d-%d", *homeGoals, *awayGoals)
		}
		fmt.Printf("   Match %d: %s %s %s\n", matchID, homeTeam, score, awayTeam)
	}
	return nil
}