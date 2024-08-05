package main

import (
	"encoding/csv"
	"log"
	"os"
	"github.com/yourusername/blue-jay/config" 
	"github.com/gocolly/colly"
)

type Game struct {
	Title       string
	Description string
	Price       string
}

func scrapeGameList(cfg *config.Config) []Game {
	var games []Game

	collector := colly.NewCollector(
		colly.UserAgent(cfg.UserAgent),
	)

	collector.SetRequestTimeout(cfg.RequestTimeout)

	collector.OnHTML(".search-results-row", func(e *colly.HTMLElement) {
		game := Game{
			Title:       e.ChildText(".search-results-row-game-title"),
			Description: e.ChildText(".search-results-row-game-infos"),
			Price:       e.ChildText(".search-results-row-price"),
		}
		games = append(games, game)
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with response: %v\n", r.Request.URL, err)
	})

	err := collector.Visit(cfg.GameListURL)
	if err != nil {
		log.Fatalf("Failed to visit URL %s: %v", cfg.GameListURL, err)
	}

	collector.Wait()

	return games
}

func writeToCSV(games []Game, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{"Title", "Description", "Price"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, game := range games {
		record := []string{game.Title, game.Description, game.Price}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	writer.Flush()

	return writer.Error()
}

func main() {
	config.CheckEnvironmentVariables()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	games := scrapeGameList(cfg)

	err = writeToCSV(games, cfg.CSVFileName)
	if err != nil {
		log.Fatalf("Failed to write games to CSV: %v", err)
	}
}

