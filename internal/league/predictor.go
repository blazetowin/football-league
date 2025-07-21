package league


import (
	"fmt"
	"math"
	"sort"
	"go-football-league/internal/models"
)

// Şampiyonluk oranlarını hesapla ve yazdır
func PrintChampionshipPredictions(week int, table []models.LeagueTableRow) {
	if week < 4 {
		return // 4. haftadan önce gösterme
	}

	fmt.Printf("\n🏆 %d. Hafta Şampiyonluk Tahminleri:\n", week)

	// 1. Toplam puanı topla (normalize için)
	totalPoints := 0
	for _, t := range table {
		totalPoints += t.Points
	}

	// 2. Boş puan varsa rastgele dağılmasın
	if totalPoints == 0 {
		for _, t := range table {
			fmt.Printf("%-15s %%25\n", t.TeamName)
		}
		return
	}

	// 3. Oran hesapla
	type Prediction struct {
		Team  string
		Oran  float64
	}
	var preds []Prediction
	for _, t := range table {
		p := float64(t.Points) / float64(totalPoints) * 100
		preds = append(preds, Prediction{
			Team: t.TeamName,
			Oran: p,
		})
	}

	// 4. Büyükten küçüğe sırala
	sort.Slice(preds, func(i, j int) bool {
		return preds[i].Oran > preds[j].Oran
	})

	// 5. Yazdır
	for _, p := range preds {
		// En fazla 1 ondalık göster
		fmt.Printf("%-15s %%%.0f\n", p.Team, math.Round(p.Oran))
	}
}
