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

// ✅ Ana router fonksiyonu
func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/matches/{week}", GetWeekMatches).Methods("GET")
	r.HandleFunc("/api/league-table", GetLeagueTable).Methods("GET")
	r.HandleFunc("/api/match/{id}", UpdateMatchScore).Methods("PUT")
	r.HandleFunc("/api/play-all-weeks", PlayAllWeeks).Methods("GET")
	r.HandleFunc("/api/week-summary", GetWeekSummary).Methods("GET")
	r.HandleFunc("/api/championship-predictions/{week}", GetChampionshipPredictions).Methods("GET") // 🔧 Burada düzeltildi

	return r
}

// ✅ Belirli haftanın maçlarını oluşturur, simüle eder ve döner
func GetWeekMatches(w http.ResponseWriter, r *http.Request) {
	weekStr := mux.Vars(r)["week"]
	week, err := strconv.Atoi(weekStr)
	if err != nil {
		http.Error(w, "Geçersiz hafta numarası", http.StatusBadRequest)
		return
	}

	err = league.GenerateWeeklyMatches(week)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = league.SimulateScores(week)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	matches, err := league.GetMatchesByWeek(week)
	if err != nil {
		http.Error(w, "Maçlar alınamadı", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}

// ✅ Lig puan tablosunu döner
func GetLeagueTable(w http.ResponseWriter, r *http.Request) {
	weekStr := r.URL.Query().Get("week")
	week, err := strconv.Atoi(weekStr)
	if err != nil || week < 1 || week > 6 {
		http.Error(w, "Geçersiz veya eksik 'week' parametresi", http.StatusBadRequest)
		return
	}

	table, err := league.GenerateLeagueTable(week)
	if err != nil {
		http.Error(w, "Puan durumu getirilemedi", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(table)
}

// ✅ Maç skorunu günceller
func UpdateMatchScore(w http.ResponseWriter, r *http.Request) {
	matchIDStr := mux.Vars(r)["id"]
	matchID, err := strconv.Atoi(matchIDStr)
	if err != nil {
		http.Error(w, "Geçersiz maç ID", http.StatusBadRequest)
		return
	}

	var update struct {
		HomeGoals int `json:"home_goals"`
		AwayGoals int `json:"away_goals"`
	}

	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "İstek verisi okunamadı", http.StatusBadRequest)
		return
	}

	err = league.UpdateMatchResult(matchID, update.HomeGoals, update.AwayGoals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("✅ Maç skoru güncellendi"))
}

// ✅ Tüm haftaları otomatik oynatır
func PlayAllWeeks(w http.ResponseWriter, r *http.Request) {
	results := make(map[int]interface{})
	for week := 1; week <= 6; week++ {
		err := league.GenerateWeeklyMatches(week)
		if err != nil {
			http.Error(w, fmt.Sprintf("Week %d fixture error: %v", week, err), http.StatusInternalServerError)
			return
		}
		err = league.SimulateScores(week)
		if err != nil {
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// ✅ Championship Predictions endpoint
func GetChampionshipPredictions(w http.ResponseWriter, r *http.Request) {
	weekStr := mux.Vars(r)["week"]
	week, err := strconv.Atoi(weekStr)
	if err != nil {
		http.Error(w, "Geçersiz hafta", http.StatusBadRequest)
		return
	}
	predictions, err := generateChampionshipPredictions(week)
	if err != nil {
		http.Error(w, "Tahminler alınamadı", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(predictions)
}

// ✅ Tahmin hesaplama fonksiyonu
func generateChampionshipPredictions(week int) ([]map[string]interface{}, error) {
	if week < 4 {
		return []map[string]interface{}{}, nil
	}
	table, err := league.GenerateLeagueTable(week)
	if err != nil {
		return nil, err
	}
	totalPoints := 0
	for _, t := range table {
		totalPoints += t.Points
	}
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
	sort.Slice(preds, func(i, j int) bool {
		return preds[i].Rate > preds[j].Rate
	})
	var response []map[string]interface{}
	for _, p := range preds {
		response = append(response, map[string]interface{}{
			"team":   p.Team,
			"chance": int(p.Rate),
		})
	}
	return response, nil
}

// ✅ Hafta özeti
func GetWeekSummary(w http.ResponseWriter, r *http.Request) {
	weekStr := r.URL.Query().Get("week")
	week, err := strconv.Atoi(weekStr)
	if err != nil || week < 1 || week > 6 {
		http.Error(w, "Geçersiz hafta", http.StatusBadRequest)
		return
	}
	matches, err := league.GetMatchesByWeek(week)
	if err != nil {
		http.Error(w, "Maçlar alınamadı", http.StatusInternalServerError)
		return
	}
	table, err := league.GenerateLeagueTable(week)
	if err != nil {
		http.Error(w, "Puan durumu alınamadı", http.StatusInternalServerError)
		return
	}
	predictions, err := generateChampionshipPredictions(week)
	if err != nil {
		http.Error(w, "Tahminler hesaplanamadı", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"week":        week,
		"matches":     matches,
		"leagueTable": table,
		"predictions": predictions,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
