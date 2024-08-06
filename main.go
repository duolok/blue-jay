package main

import (
	"fmt"
	"os"
	"log"

	"github.com/duolok/blue-jay/scrapers/instant_gaming"
)

func loadScrapers() []string {
	var scrapers []string

	items, _ := os.ReadDir("./scrapers/")
	for _, item := range items {
		if item.IsDir() {
			scrapers = append(scrapers, item.Name())
		}
	}

	return scrapers
}


func main() {
	scrapers := loadScrapers()
	for _, item := range scrapers {
		fmt.Println(item)
	}


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

