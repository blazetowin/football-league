package league

import (
	"fmt" // Importing fmt for formatted I/O
	"sort" 	// Importing sort for sorting slices
	"go-football-league/internal/models" // Importing models for LeagueTableRow type
	"go-football-league/internal/storage" 	// Importing storage for database operations
)

// GenerateLeagueTable computes the league standings from all matches
// It aggregates match results to calculate points, wins, losses, and goal differences for each team
// It returns a slice of LeagueTableRow containing the standings
// If an error occurs during database operations, it returns the error
func GenerateLeagueTable() ([]models.LeagueTableRow, error) {
	rows, err := storage.DB.Query(`
		SELECT 
			m.home_team_id, t1.name, m.home_goals, m.away_goals,
			m.away_team_id, t2.name
		FROM matches m
		JOIN teams t1 ON m.home_team_id = t1.id
		JOIN teams t2 ON m.away_team_id = t2.id
		WHERE m.home_goals IS NOT NULL AND m.away_goals IS NOT NULL
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed after use

	// Initialize a map to hold team statistics
	// This map will store LeagueTableRow for each team
	// The key is the team ID and the value is a pointer to LeagueTableRow
	// This allows us to easily update the statistics for each team as we process matches
	stats := make(map[int]*models.LeagueTableRow) 

	// Iterate over the result set
	// For each match we will update the statistics for both home and away teams
	for rows.Next() {
		var homeID, awayID, homeGoals, awayGoals int
		var homeName, awayName string

		if err := rows.Scan(&homeID, &homeName, &homeGoals, &awayGoals, &awayID, &awayName); err != nil {
			return nil, err
		}

		// Initialize rows if needed
		// If the team is not already in the stats map we create a new LeagueTableRow
		// This ensures that we have a row for each team in the league
		// We use the team ID as the key and the team name for display purposes
		// This allows us to keep track of each team's performance throughout the season
		if _, ok := stats[homeID]; !ok {
			stats[homeID] = &models.LeagueTableRow{TeamID: homeID, TeamName: homeName}
		}
		if _, ok := stats[awayID]; !ok {
			stats[awayID] = &models.LeagueTableRow{TeamID: awayID, TeamName: awayName}
		}
		// Update statistics for home team
		// We increment the played matches count for both teams
		// Depending on the match result we update wins, losses, and draws
		// The home team gets 3 points for a win, 1 point for a draw, and 0 points for a loss
		// The away team gets the same points
		// This allows us to calculate the league standings based on match results
		home := stats[homeID]
		away := stats[awayID]

		home.Played++
		away.Played++	

		home.GoalsFor += homeGoals
		home.GoalsAgainst += awayGoals

		away.GoalsFor += awayGoals
		away.GoalsAgainst += homeGoals

		if homeGoals > awayGoals {
			home.Wins++
			home.Points += 3
			away.Losses++
		} else if awayGoals > homeGoals {
			away.Wins++
			away.Points += 3
			home.Losses++
		} else {
			home.Draws++
			away.Draws++
			home.Points++
			away.Points++
		}
	}
	// After processing all matches we need to calculate the goal difference for each team
	// The goal difference is calculated as Goals For - Goals Against
	// This is done for each team in the stats map
	var table []models.LeagueTableRow
	for _, row := range stats {
		row.GoalDiff = row.GoalsFor - row.GoalsAgainst
		table = append(table, *row)
	}
	// We sort the league table based on points, goal difference, and goals for
	// This ensures that the team with the most points is at the top of the table
	// If two teams have the same points, the one with the better goal difference comes first
	// If they also have the same goal difference, the one with more goals scored is ranked
	sort.SliceStable(table, func(i, j int) bool {
		if table[i].Points != table[j].Points {
			return table[i].Points > table[j].Points
		}
		if table[i].GoalDiff != table[j].GoalDiff {
			return table[i].GoalDiff > table[j].GoalDiff
		}
		return table[i].GoalsFor > table[j].GoalsFor
	})
	// This is the final result of the league standings after processing all matches
	// The table contains LeagueTableRow for each team with their statistics
	fmt.Println("📊 League table generated.")
	return table, nil
}
