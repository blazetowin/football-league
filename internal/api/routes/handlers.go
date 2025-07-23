package routes

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
	"go-football-league/internal/league"
)

// SetupRouter initializes and returns the main API router with all endpoints registered.
func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Registering HTTP route handlers
	r.HandleFunc("/api/matches/{week}", GetWeekMatches).Methods("GET")
	r.HandleFunc("/api/league-table", GetLeagueTable).Methods("GET")
	r.HandleFunc("/api/match/{id}", UpdateMatchScore).Methods("PUT")
	r.HandleFunc("/api/play-all-weeks", PlayAllWeeks).Methods("GET")
	r.HandleFunc("/api/week-summary", GetWeekSummary).Methods("GET")
	r.HandleFunc("/api/championship-predictions/{week}", GetChampionshipPredictions).Methods("GET")

	return r
}

// GetWeekMatches handles GET /api/matches/{week}
// It generates fixtures, simulates scores, and returns all matches for the given week.
func GetWeekMatches(w http.ResponseWriter, r *http.Request) {
	weekStr := mux.Vars(r)["week"]
	week, err := strconv.Atoi(weekStr)
	if err != nil {
		http.Error(w, "Invalid week number", http.StatusBadRequest)
		return
	}

	// Generate match fixtures for the specified week
	if err := league.GenerateWeeklyMatches(week); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Simulate match scores
	if err := league.SimulateScores(week); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch all matches for the given week
	matches, err := league.GetMatchesByWeek(week)
	if err != nil {
		http.Error(w, "Failed to retrieve matches", http.StatusInternalServerError)
		return
	}

	// Return matches as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}

// GetLeagueTable handles GET /api/league-table?week=
// Returns the league standings for a given week.
func GetLeagueTable(w http.ResponseWriter, r *http.Request) {
	weekStr := r.URL.Query().Get("week")
	week, err := strconv.Atoi(weekStr)
	if err != nil || week < 1 || week > 6 {
		http.Error(w, "Invalid or missing 'week' parameter", http.StatusBadRequest)
		return
	}

	// Generate the league standings
	table, err := league.GenerateLeagueTable(week)
	if err != nil {
		http.Error(w, "Failed to generate league table", http.StatusInternalServerError)
		return
	}

	// Return standings as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(table)
}

// UpdateMatchScore handles PUT /api/match/{id}
// Updates the match result with provided home and away goals.
func UpdateMatchScore(w http.ResponseWriter, r *http.Request) {
	matchIDStr := mux.Vars(r)["id"]
	matchID, err := strconv.Atoi(matchIDStr)
	if err != nil {
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	// Parse request body to extract new score
	var update struct {
		HomeGoals int `json:"home_goals"`
		AwayGoals int `json:"away_goals"`
	}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Apply the score update
	if err := league.UpdateMatchResult(matchID, update.HomeGoals, update.AwayGoals); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Match score updated successfully"))
}

// PlayAllWeeks handles GET /api/play-all-weeks
// Simulates all league weeks (1 to 6) and returns the results for each week.
func PlayAllWeeks(w http.ResponseWriter, r *http.Request) {
	results := make(map[int]interface{})
	for week := 1; week <= 6; week++ {
		if err := league.GenerateWeeklyMatches(week); err != nil {
			http.Error(w, fmt.Sprintf("Week %d fixture error: %v", week, err), http.StatusInternalServerError)
			return
		}
		if err := league.SimulateScores(week); err != nil {
			http.Error(w, fmt.Sprintf("Week %d simulation error: %v", week, err), http.StatusInternalServerError)
			return
		}
		matches, err := league.GetMatchesByWeek(week)
		if err != nil {
			http.Error(w, fmt.Sprintf("Week %d matches fetch error: %v", week, err), http.StatusInternalServerError)
			return
		}
		results[week] = matches
	}

	// Return match results of all weeks
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// GetChampionshipPredictions handles GET /api/championship-predictions/{week}
// Calculates title-winning probabilities based on current standings.
func GetChampionshipPredictions(w http.ResponseWriter, r *http.Request) {
	weekStr := mux.Vars(r)["week"]
	week, err := strconv.Atoi(weekStr)
	if err != nil {
		http.Error(w, "Invalid week", http.StatusBadRequest)
		return
	}

	predictions, err := generateChampionshipPredictions(week)
	if err != nil {
		http.Error(w, "Failed to compute predictions", http.StatusInternalServerError)
		return
	}

	// Return probability list as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(predictions)
}

// generateChampionshipPredictions computes winning probability of each team
// based on their points in the standings as of the given week.
// Returns nil if the simulation is called before week 4.
func generateChampionshipPredictions(week int) ([]map[string]interface{}, error) {
	if week < 4 {
		// Not enough data to predict before week 4
		return []map[string]interface{}{}, nil
	}

	// Fetch league table
	table, err := league.GenerateLeagueTable(week)
	if err != nil {
		return nil, err
	}

	// Calculate total points to determine percentage weights
	totalPoints := 0
	for _, t := range table {
		totalPoints += t.Points
	}

	// Store team and their winning chance
	type Prediction struct {
		Team string
		Rate float64
	}
	var preds []Prediction
	for _, t := range table {
		var rate float64
		if totalPoints > 0 {
			rate = (float64(t.Points) / float64(totalPoints)) * 100
		}
		preds = append(preds, Prediction{
			Team: t.TeamName,
			Rate: math.Round(rate),
		})
	}

	// Sort by probability descending
	sort.Slice(preds, func(i, j int) bool {
		return preds[i].Rate > preds[j].Rate
	})

	// Format response
	var response []map[string]interface{}
	for _, p := range preds {
		response = append(response, map[string]interface{}{
			"team":   p.Team,
			"chance": int(p.Rate),
		})
	}
	return response, nil
}

// GetWeekSummary handles GET /api/week-summary?week=
// Returns a weekly summary including matches, league table, and predictions.
func GetWeekSummary(w http.ResponseWriter, r *http.Request) {
	weekStr := r.URL.Query().Get("week")
	week, err := strconv.Atoi(weekStr)
	if err != nil || week < 1 || week > 6 {
		http.Error(w, "Invalid week", http.StatusBadRequest)
		return
	}

	// Fetch all required data
	matches, err := league.GetMatchesByWeek(week)
	if err != nil {
		http.Error(w, "Failed to fetch matches", http.StatusInternalServerError)
		return
	}
	table, err := league.GenerateLeagueTable(week)
	if err != nil {
		http.Error(w, "Failed to fetch league table", http.StatusInternalServerError)
		return
	}
	predictions, err := generateChampionshipPredictions(week)
	if err != nil {
		http.Error(w, "Failed to fetch predictions", http.StatusInternalServerError)
		return
	}

	// Aggregate and return full weekly summary
	response := map[string]interface{}{
		"week":        week,
		"matches":     matches,
		"leagueTable": table,
		"predictions": predictions,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
