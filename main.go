package main

import (
	"fmt"
	"os"
	"log"
	"bufio"

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


func loadLastSearch(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text()) 
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return lines, nil
}

func main() {
	scrapers := loadScrapers()
	for _, item := range scrapers {
		fmt.Println(item)
	}

	lines, err := loadLastSearch("games.csv")
	if err != nil {
		fmt.Println("Error: ", err)
		return 
	}

	for _, line := range(lines) {
		fmt.Println(line)
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

