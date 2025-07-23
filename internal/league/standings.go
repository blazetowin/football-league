package league

import (
	"fmt"
	"sort"

	models "go-football-league/internal/domain"
	storage "go-football-league/internal/repository"
)

// GenerateLeagueTable computes the league standings.
// It reads played matches from the database and calculates total points, goals, wins, losses and draws for each team. The final table is sorted by points, gd and goals scored.
func GenerateLeagueTable(upToWeek int) ([]models.LeagueTableRow, error) {
	// Query all played matches up to the specified week
	rows, err := storage.DB.Query(`
		SELECT 
			m.home_team_id, t1.name, m.home_goals, m.away_goals,
			m.away_team_id, t2.name
		FROM matches m
		JOIN teams t1 ON m.home_team_id = t1.id
		JOIN teams t2 ON m.away_team_id = t2.id
		WHERE m.week <= ? AND m.home_goals IS NOT NULL AND m.away_goals IS NOT NULL
	`, upToWeek)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Use a map to collect statistics per team
	stats := make(map[int]*models.LeagueTableRow)

	for rows.Next() {
		var homeID, awayID, homeGoals, awayGoals int
		var homeName, awayName string

		// Read one match result
		if err := rows.Scan(&homeID, &homeName, &homeGoals, &awayGoals, &awayID, &awayName); err != nil {
			return nil, err
		}

		// Initialize team rows if not already added
		if _, ok := stats[homeID]; !ok {
			stats[homeID] = &models.LeagueTableRow{TeamID: homeID, TeamName: homeName}
		}
		if _, ok := stats[awayID]; !ok {
			stats[awayID] = &models.LeagueTableRow{TeamID: awayID, TeamName: awayName}
		}

		home := stats[homeID]
		away := stats[awayID]

		// Update matches played
		home.Played++
		away.Played++

		// Update goals scored and conceded
		home.GoalsFor += homeGoals
		home.GoalsAgainst += awayGoals
		away.GoalsFor += awayGoals
		away.GoalsAgainst += homeGoals

		// Assign points and match results
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

	// Finalize table with goal difference and flatten the map
	var table []models.LeagueTableRow
	for _, row := range stats {
		row.GoalDiff = row.GoalsFor - row.GoalsAgainst
		table = append(table, *row)
	}

	// Sort the table: Points > Goal Difference > Goals For
	sort.SliceStable(table, func(i, j int) bool {
		if table[i].Points != table[j].Points {
			return table[i].Points > table[j].Points
		}
		if table[i].GoalDiff != table[j].GoalDiff {
			return table[i].GoalDiff > table[j].GoalDiff
		}
		return table[i].GoalsFor > table[j].GoalsFor
	})

	fmt.Println("League table generated.")
	return table, nil
}
