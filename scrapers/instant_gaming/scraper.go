package instant_gaming

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Game struct {
	Title string
	Link  string
	Price string
}

type InstantGamingScraper struct {
	Config *Config
	Games  []Game
}

func NewInstantGamingScraper(cfg *Config) *InstantGamingScraper {
	return &InstantGamingScraper{
		Config: cfg,
	}
}

func (s *InstantGamingScraper) Scrape(gameName string) {
	collector := colly.NewCollector(
		colly.UserAgent(s.Config.UserAgent),
		colly.Async(true),
	)

	collector.SetRequestTimeout(time.Duration(s.Config.RequestTimeout) * time.Second)
	collector.OnHTML(".search .item", func(h *colly.HTMLElement) {
		game := Game{
			Title: h.ChildText(".title"),
			Link:  h.ChildAttr("a", "href"),
			Price: h.ChildText(".price"),
		}

		if game.Title != "" {
			s.Games = append(s.Games, game)
		}
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with response: %v\n", r.Request.URL, err)
	})

	err := collector.Visit(s.Config.GameListURL + transformGameSearchString(gameName))
	if err != nil {
		log.Fatalf("Failed to visit URL %s: %v", s.Config.GameListURL, err)
	}

	collector.Wait()
}

func (s *InstantGamingScraper) WriteToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{"Title", "Price", "Link"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, game := range s.Games {
		record := []string{game.Title, game.Price, game.Link}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
func (s *InstantGamingScraper) GetPrices() {
	// Implement price fetching logic if needed
}

func transformGameSearchString(input string) string {
	return strings.ReplaceAll(input, " ", "%20")
}

