package models

// Package models defines the data structures used across the football league simulation.
// It includes core types such as teams, matches, league standings, and prediction models.

// Team represents a football team with a unique ID, name, and power rating used for match simulations.
type Team struct {
	ID    int
	Name  string
	Power int
}

// Match represents a football match between two teams during a specific week.
// It also includes simulated or updated scores and team names for display purposes.
type Match struct {
	ID           int    
	Week         int    
	HomeTeamID   int    
	AwayTeamID   int    
	HomeGoals    int
	AwayGoals    int
	HomeTeamName string
	AwayTeamName string
}

// LeagueTableRow represents the position and performance statistics of a team in the league standings.
type LeagueTableRow struct {
	TeamID       int 
	TeamName     string 
	Points       int 
	Played       int 
	Wins         int 
	Draws        int 
	Losses       int 
	GoalsFor     int 
	GoalsAgainst int 
	GoalDiff     int 
}

// Prediction represents a team's probability of winning the league championship based on current standings.
type Prediction struct {
	TeamID   int     
	TeamName string  
	Chance   float64 
}
