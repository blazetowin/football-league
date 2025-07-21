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
func GenerateLeagueTable(upToWeek int) ([]models.LeagueTableRow, error) {
	rows, err := storage.DB.Query(`
		SELECT 
			m.home_team_id, t1.name, m.home_goals, m.away_goals,
			m.away_team_id, t2.name
		FROM matches m
		JOIN teams t1 ON m.home_team_id = t1.id
		JOIN teams t2 ON m.away_team_id = t2.id
		WHERE m.week <= ? AND m.home_goals IS NOT NULL AND m.away_goals IS NOT NULL
	`, upToWeek) // 👈 Sadece belirli haftaya kadar olan maçlar
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[int]*models.LeagueTableRow)

	for rows.Next() {
		var homeID, awayID, homeGoals, awayGoals int
		var homeName, awayName string

		if err := rows.Scan(&homeID, &homeName, &homeGoals, &awayGoals, &awayID, &awayName); err != nil {
			return nil, err
		}

		if _, ok := stats[homeID]; !ok {
			stats[homeID] = &models.LeagueTableRow{TeamID: homeID, TeamName: homeName}
		}
		if _, ok := stats[awayID]; !ok {
			stats[awayID] = &models.LeagueTableRow{TeamID: awayID, TeamName: awayName}
		}

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

	var table []models.LeagueTableRow
	for _, row := range stats {
		row.GoalDiff = row.GoalsFor - row.GoalsAgainst
		table = append(table, *row)
	}

	sort.SliceStable(table, func(i, j int) bool {
		if table[i].Points != table[j].Points {
			return table[i].Points > table[j].Points
		}
		if table[i].GoalDiff != table[j].GoalDiff {
			return table[i].GoalDiff > table[j].GoalDiff
		}
		return table[i].GoalsFor > table[j].GoalsFor
	})

	fmt.Println("📊 League table generated.")
	return table, nil
}
