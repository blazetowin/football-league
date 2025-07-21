package league

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"go-football-league/internal/models"
	"go-football-league/internal/storage"
)

// Haftalık maçları oluşturur
func GenerateWeeklyMatches(week int) error {
	var count int
	err := storage.DB.QueryRow("SELECT COUNT(*) FROM matches WHERE week = ?", week).Scan(&count)
	if err != nil {
		return err
	}
	if count == 2 {
		fmt.Printf("\u2139\ufe0f Matches already exist for week %d\n", week)
		return nil
	}
	if count == 0 {
		return errors.New("\u274c Fixture not created — please run CreateFixture() first.")
	}
	if count != 2 {
		return fmt.Errorf("\u274c Unexpected match count for week %d: expected 2, got %d", week, count)
	}
	return nil
}

// Maçlara skor simülasyonu uygular
func SimulateScores(week int) error {
	rows, err := storage.DB.Query(`
		SELECT m.id, t1.power, t2.power
		FROM matches m
		JOIN teams t1 ON m.home_team_id = t1.id
		JOIN teams t2 ON m.away_team_id = t2.id
		WHERE m.week = ? AND m.home_goals IS NULL AND m.away_goals IS NULL
	`, week)
	if err != nil {
		return err
	}
	defer rows.Close()
	type match struct {
		ID        int
		PowerHome int
		PowerAway int
	}
	var matches []match
	for rows.Next() {
		var m match
		if err := rows.Scan(&m.ID, &m.PowerHome, &m.PowerAway); err != nil {
			return err
		}
		matches = append(matches, m)
	}

	rand.Seed(time.Now().UnixNano())
	for _, m := range matches {
		homeGoals := rand.Intn(min((m.PowerHome/10)+2+1, 6))
		awayGoals := rand.Intn(min((m.PowerAway/10)+2, 6))

		fmt.Printf(" Match %d simulated → Home: %d | Away: %d\n", m.ID, homeGoals, awayGoals)
		res, err := storage.DB.Exec(`
			UPDATE matches SET home_goals = ?, away_goals = ? WHERE id = ?
		`, homeGoals, awayGoals, m.ID)
		if err != nil {
			fmt.Println("\u274c UPDATE error for match", m.ID, ":", err)
			return err
		}
		rowsAffected, _ := res.RowsAffected()
		fmt.Printf(" Match %d update affected rows: %d\n", m.ID, rowsAffected)
	}
	return nil
}

// Fixture oluşturur
func CreateFixture() error {
	var existing int
	err := storage.DB.QueryRow("SELECT COUNT(*) FROM matches").Scan(&existing)
	if err != nil {
		return fmt.Errorf("Fixture kontrolü başarısız: %v", err)
	}
	if existing > 0 {
		fmt.Println("\u2139\ufe0f Fixture zaten oluşturulmuş. Yeniden oluşturulmuyor.")
		return nil
	}

	rows, err := storage.DB.Query("SELECT id FROM teams ORDER BY id")
	if err != nil {
		return err
	}
	defer rows.Close()

	var teamIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		teamIDs = append(teamIDs, id)
	}

	if len(teamIDs) != 4 {
		return errors.New("Fixture sadece 4 takım içindir")
	}

	type Match struct {
		Week int
		Home int
		Away int
	}
	fixture := []Match{
		{1, teamIDs[0], teamIDs[1]}, {1, teamIDs[2], teamIDs[3]},
		{2, teamIDs[0], teamIDs[2]}, {2, teamIDs[1], teamIDs[3]},
		{3, teamIDs[0], teamIDs[3]}, {3, teamIDs[1], teamIDs[2]},
		{4, teamIDs[1], teamIDs[0]}, {4, teamIDs[3], teamIDs[2]},
		{5, teamIDs[2], teamIDs[0]}, {5, teamIDs[3], teamIDs[1]},
		{6, teamIDs[3], teamIDs[0]}, {6, teamIDs[2], teamIDs[1]},
	}

	for _, match := range fixture {
		_, err := storage.DB.Exec(`
			INSERT INTO matches (week, home_team_id, away_team_id, home_goals, away_goals)
			VALUES (?, ?, ?, NULL, NULL)
		`, match.Week, match.Home, match.Away)
		if err != nil {
			return fmt.Errorf("Match kaydı başarısız: %v", err)
		}
	}

	fmt.Println("\u2705 Fixture created successfully.")
	return nil
}

// Belirli haftaya ait maçları getirir
func GetMatchesByWeek(week int) ([]models.Match, error) {
	rows, err := storage.DB.Query(`
		SELECT m.id, m.week, m.home_team_id, m.away_team_id, m.home_goals, m.away_goals,
		       ht.name as home_team_name, at.name as away_team_name
		FROM matches m
		JOIN teams ht ON m.home_team_id = ht.id
		JOIN teams at ON m.away_team_id = at.id
		WHERE m.week = ?
	`, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []models.Match
	for rows.Next() {
		var m models.Match
		err := rows.Scan(&m.ID, &m.Week, &m.HomeTeamID, &m.AwayTeamID, &m.HomeGoals, &m.AwayGoals, &m.HomeTeamName, &m.AwayTeamName)
		if err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}

	return matches, nil
}

// Maç sonucunu günceller
func UpdateMatchResult(matchID int, homeGoals, awayGoals int) error {
	_, err := storage.DB.Exec(`
		UPDATE matches
		SET home_goals = ?, away_goals = ?
		WHERE id = ?
	`, homeGoals, awayGoals, matchID)
	return err
}

// Yardımcı fonksiyon: min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
