package interfaces

type Scraper interface { 
    scrape(gameName string)
    writeToFile(path string)
    getPrices()
}
