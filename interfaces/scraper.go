package interfaces

type Scraper interface { 
    Scrape(gameName string)
    WriteToFile(path string) error
    GetPrices() []string
}
