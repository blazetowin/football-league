package models // Package models defines the data structures used in the football league application.
// It includes types for teams, matches, league table rows, and predictions.

type Team struct {
	ID    int
	Name  string
	Power int
}

type Match struct {
	ID         int
	Week       int
	HomeTeamID int
	AwayTeamID int
	HomeGoals  int
	AwayGoals  int
	HomeTeamName string
	AwayTeamName string
}

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

type Prediction struct {
	TeamID   int
	TeamName string
	Chance   float64
}
