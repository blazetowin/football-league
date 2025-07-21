package main

import (
	"fmt"
	"log"

	"go-football-league/internal/league"
	"go-football-league/internal/storage"
)

func main() {
	// Veritabanı bağlantısını başlat
	storage.Connect()


	// 4 hafta boyunca lig oynatılacak
	for week := 1; week <= 4; week++ {
		fmt.Printf("\n🗓️ ====================== %d. HAFTA ======================\n\n", week)

		// Hafta oynatılıyor (eğer daha önce oynanmadıysa)
		err := league.PlayWeek(week)
		if err != nil {
			log.Fatalf("❌ Hafta %d oynatılamadı: %v", week, err)
		}

		// Maç sonuçlarını yazdır
		fmt.Printf("\n🏟️ %d. Hafta Maç Sonuçları:\n", week)
		err = league.PrintMatchesOfWeek(week)
		if err != nil {
			log.Fatalf("❌ Maçlar listelenemedi: %v", err)
		}

		// Güncel puan durumunu yazdır
		fmt.Printf("\n📊 %d. Hafta Sonu Puan Durumu:\n", week)
		table, err := league.GenerateLeagueTable()
		if err != nil {
			log.Fatalf("❌ Puan durumu getirilemedi: %v", err)
		}

		fmt.Println("------------------------------------------------------------")
		fmt.Printf("%-15s %2s %2s %2s %2s %4s %4s %4s %4s\n", "Takım", "O", "G", "B", "M", "AG", "YG", "+/-", "P")
		fmt.Println("------------------------------------------------------------")
		for _, row := range table {
			fmt.Printf("%-15s %2d %2d %2d %2d %4d %4d %4d %4d\n",
				row.TeamName, row.Played, row.Wins, row.Draws, row.Losses,
				row.GoalsFor, row.GoalsAgainst, row.GoalDiff, row.Points)
		}
		fmt.Println("------------------------------------------------------------")
	}
}
