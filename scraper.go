package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	"time"
	"github.com/duolok/blue-jay/config"
	"github.com/gocolly/colly"
)

type Game struct {
	Title string
	Link  string
	Price string
}

func scrapeGameList(cfg *config.Config, searchedGame string) []Game {
	var games []Game

	collector := colly.NewCollector(
		colly.UserAgent(cfg.UserAgent),
		colly.Async(true),
	)

	collector.SetRequestTimeout(time.Duration(cfg.RequestTimeout) * time.Second)
	collector.OnHTML(".search .item", func(h *colly.HTMLElement) {
		game := Game{
			Title: h.ChildText(".title"),
			Link: h.ChildAttr("a", "href"),
			Price: h.ChildText(".price"),
		}

		if game.Title != "" {
			games = append(games, game)
		}
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with response: %v\n", r.Request.URL, err)
	})

	err := collector.Visit(cfg.GameListURL + transformGameSearchString(searchedGame))
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

	headers := []string{"Title", "Price"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, game := range games {
		record := []string{game.Title, game.Price, game.Link}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func transformGameSearchString(input string) string {
	converted := strings.ReplaceAll(input, " ", "%20")
	return converted + "/"
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	games := scrapeGameList(cfg, "Elden Ring")

	err = writeToCSV(games, cfg.CSVFileName)
	if err != nil {
		log.Fatalf("Failed to write games to CSV: %v", err)
	}
}

