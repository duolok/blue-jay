package main

import (
	"log"

	"github.com/duolok/blue-jay/scrapers/instant_gaming"
)

func main() {
	cfg, err := instant_gaming.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	scraper := instant_gaming.NewInstantGamingScraper(cfg)
	scraper.Scrape("Hollow Knight")

	err = scraper.WriteToFile(cfg.CSVFileName)
	if err != nil {
		log.Fatalf("Failed to write games to CSV: %v", err)
	}
}

