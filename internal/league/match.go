package league // Package league handles the generation and simulation of matches in a football league.

import (
	"errors"       // Importing errors package to define custom error messages
	"fmt"          // Importing fmt for logging
	"math/rand"    // Importing math/rand for random number generation
	"time"         // Importing time for time-related functions

	"go-football-league/internal/storage" // Importing storage for database operations
)

// GenerateWeeklyMatches creates matches for a given week (if not already created)
func GenerateWeeklyMatches(week int) error {
	// Fetch all teams
	rows, err := storage.DB.Query("SELECT id FROM teams") // Query to get all team IDs
	if err != nil {  		
		return err 
	}
	defer rows.Close() // Ensure rows are closed after use

	// Collect team IDs
	// This will be used to create matches
	// We will randomly pair teams for matches
	var teamIDs []int
	for rows.Next() {
		// Scan each team ID into the teamIDs slice
		var id int 
		if err := rows.Scan(&id); err != nil { 
			return err
		}
		// Append the team ID to the slice
		teamIDs = append(teamIDs, id)
	}

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	// Shuffle the team IDs slice to randomize matchups
	// This ensures that teams are paired randomly each week
	// rand.Shuffle is a function that randomly shuffles the elements of a slice
	// It takes the length of the slice and a function that swaps two elements
	// Here we swap elements at indices i and j this creates a random order of teams for match generation
	rand.Shuffle(len(teamIDs), func(i, j int) {
		teamIDs[i], teamIDs[j] = teamIDs[j], teamIDs[i]	
	})

	// Check if matches for this week already exist,
	// if so, we don't need to generate them again
	// This prevents duplicate matches from being created
	var count int
	err = storage.DB.QueryRow("SELECT COUNT(*) FROM matches WHERE week = ?", week).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	if count >= 2 {
		return errors.New("matches for this week already exist")
	}

	// If no matches exist for this week, we proceed to create them
	// We will create matches in pairs of two teams per match
	// This means we will have len(teamIDs)/2 matches for this week
	if len(teamIDs)%2 != 0 {
		return errors.New("cannot generate matches with an odd number of teams") // Custom error message
	}

	// Generate matches (2 matches per week)
	for i := 0; i < len(teamIDs); i += 2 {
		_, err := storage.DB.Exec(`
		-- Insert a new match into the matches table
		-- This SQL command inserts a match for the given week between two teams
		-- The home_team_id and away_team_id are set to the IDs of the teams
		-- The home_goals and away_goals are initialized to NULL
		-- This means the match has not been played yet
		-- The match will be played in the future, and scores will be simulated later
			INSERT INTO matches (week, home_team_id, away_team_id, home_goals, away_goals) 
			VALUES (?, ?, ?, NULL, NULL) 	
		`, week, teamIDs[i], teamIDs[i+1])
		if err != nil {
			return err
		}
	}

	// Log the successful creation of matches for the week
	fmt.Println("ℹ️ Weekly matches inserted for week", week)
	return nil
}	

// SimulateScores generates realistic scores based on team strength
// It updates the matches for a given week with simulated scores
func SimulateScores(week int) error {
	// Select matches that haven't been simulated yet
	rows, err := storage.DB.Query(`
	-- This SQL command selects matches for a given week where the scores have not been set yet
		SELECT m.id, t1.power, t2.power
		FROM matches m
		JOIN teams t1 ON m.home_team_id = t1.id
		JOIN teams t2 ON m.away_team_id = t2.id
		WHERE m.week = ? AND m.home_goals IS NULL AND m.away_goals IS NULL
	`, week)
	if err != nil {
		return err
	}
	defer rows.Close() // Ensure rows are closed after use

	// Collect rows into a slice to close rows before updating
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

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator for score simulation
	// Iterate through the collected matches and simulate scores
	// For each match, we will calculate the scores based on team power
	// The power of each team will influence the number of goals scored
	// We will use a simple formula to generate goals based on team power and some randomness
	// The home team and the away team will have a score based on its power and some random factor

	for _, m := range matches {
		// Simulate goals for home and away teams
		// Home team gets +1 advantage
		homeGoals := rand.Intn(min((m.PowerHome/10)+2+1, 6)) // +1 for home advantage, max 5 goals, 6 possible
		// Away team gets no advantage
		// The number of goals is capped to a maximum of 5
		// This is to ensure realistic score simulation
		// The formula uses team power to determine the likelihood of scoring goals
		awayGoals := rand.Intn(min((m.PowerAway/10)+2, 6))   // max 5 goals, 6 possible

		// Log the simulated result
		fmt.Printf("🏟️ Match %d simulated → Home: %d | Away: %d\n", m.ID, homeGoals, awayGoals)

		// This SQL command updates the match with the simulated scores
		// It sets the home_goals and away_goals fields in the matches table
		// The match ID is used to identify which match to update
		// The home_goals and away_goals are set to the simulated scores
		// This means the match has been played and the scores are now recorded
		res, err := storage.DB.Exec(`
			UPDATE matches SET home_goals = ?, away_goals = ? WHERE id = ?
		`, homeGoals, awayGoals, m.ID)
		if err != nil {
			fmt.Println("❌ UPDATE error for match", m.ID, ":", err)
			return err
		}
		rowsAffected, _ := res.RowsAffected()
		fmt.Printf("🧾 Match %d update affected rows: %d\n", m.ID, rowsAffected)
	}
	return nil
}

// Helper function to cap scores to a maximum value (e.g. 5)
// This function ensures that the number of goals scored does not exceed a certain limit
// It takes two integers as input and returns the smaller of the two
// This is useful for ensuring that scores remain realistic and within a reasonable range
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
