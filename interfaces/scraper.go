package interfaces

type Scraper interface { 
    scrape(gameName string)
    writeToFile(path string) error
    getPrices() []string
}
