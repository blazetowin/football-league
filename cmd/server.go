package main

import (
	"log"
	"net/http"

	"go-football-league/internal/routes"
	"go-football-league/internal/storage"
)

func main() {
	// Veritabanı bağlantısı
	storage.Connect()

	// API route'larını başlat
	router := routes.SetupRouter()

	// HTTP sunucusunu başlat
	log.Println("🚀 Sunucu http://localhost:8080 adresinde çalışıyor...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
