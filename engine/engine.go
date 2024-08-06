package engine

import (
    "os"
    "fmt"
    "log"
    "sync"
    "bufio"
	"github.com/duolok/blue-jay/interfaces"
	"github.com/duolok/blue-jay/scrapers/instant_gaming"
)

type ScraperConstructor func(cfg *instant_gaming.Config) interfaces.Scraper

var scraperConstructors = map[string]ScraperConstructor{
	"instant_gaming": func(cfg *instant_gaming.Config) interfaces.Scraper {
		return instant_gaming.New(cfg)
	},
}

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

func search(scrapers []string, game string, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, scraperName := range scrapers {
		if constructor, exists := scraperConstructors[scraperName]; exists {
			cfg, err := instant_gaming.LoadConfig()
			if err != nil {
				log.Printf("Failed to load configuration for %s: %v", scraperName, err)
				continue
			}
			scraper := constructor(cfg)

			wg.Add(1)
			go func(scraper interfaces.Scraper, scraperName string) {
				defer wg.Done()
				scraper.Scrape(game)
				err := scraper.WriteToFile(cfg.CSVFileName)
				if err != nil {
					log.Printf("Failed to write results for %s: %v", scraperName, err)
				}
			}(scraper, scraperName)
		} else {
			log.Printf("No constructor found for scraper: %s", scraperName)
		}
	}
}
