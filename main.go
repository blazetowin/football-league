package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"go-football-league/internal/league"
	"go-football-league/internal/storage"
)

func main() {
	storage.Connect()
	fmt.Println("✅ Database connection and table setup successful.")

	// 🔁 Eğer fikstür daha önce oluşturulmamışsa oluştur
	err := league.CreateFixture()
	if err != nil {
		log.Fatalf("❌ Fixture oluşturulamadı: %v", err)
	}

	for week := 1; week <= 6; week++ {
		fmt.Printf("🗓️ ====================== %d. HAFTA ======================\n\n", week)

		err := league.GenerateWeeklyMatches(week)
		if err != nil {
			log.Fatalf("❌ Hafta %d oynatılamadı: %v", week, err)
		}

		err = league.SimulateScores(week)
		if err != nil {
			log.Fatalf("❌ Skorlar simüle edilemedi: %v", err)
		}

		fmt.Printf("\n🏟️ %d. Hafta Maç Sonuçları:\n", week)
		league.PrintMatchesOfWeek(week)

		fmt.Printf("\n📊 %d. Hafta Sonu Puan Durumu:\n", week)
		table, err := league.GenerateLeagueTable(week)
		if err != nil {
			log.Fatalf("❌ Puan durumu getirilemedi: %v", err)
		}

		league.PrintLeagueTableRows(table)
		fmt.Println()
		
		// 🔮 Şampiyonluk tahminini yazdır (4. haftadan sonra)
		league.PrintChampionshipPredictions(week, table)

		if week < 6 {
			fmt.Print("\n🔁 Devam etmek için Enter'a bas...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			fmt.Println()
		}
	}
}
