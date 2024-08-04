package main

import (
 "encoding/csv"
 "log"
 "os"
 "time"
 "github.com/gocolly/colly"
)

type gameStruct struct {
    title           string
    description     string
    price           string
}

func scrapeAndWriteCSV() []gameStruct {
 var scrapData []gameStruct

 c := colly.NewCollector(
  colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
 )

    c.SetRequestTimeout(time.Second * 10)

 c.OnHTML(".search-results-row", func(e *colly.HTMLElement) {
  game := gameStruct{}

  game.title = e.ChildText(".search-results-row-game-title")
  game.description = e.ChildText(".search-results-row-game-infos")
  game.price = e.ChildText(".search-results-row-price")

  scrapData = append(scrapData, game)
 })

 c.OnError(func(r *colly.Response, err error) {
  log.Printf("Request URL: %s failed with response: %v\n", r.Request.URL, err)
 })

 c.Visit("https://www.allkeyshop.com/blog/catalogue/search-Elden+Ring/")
 c.Wait()

 file, err := os.Create("link1.csv")
 if err != nil {
  log.Fatalln("Failed to create output CSV file", err)
 }
 defer file.Close()

 writer := csv.NewWriter(file)

 headers := []string{
  "url",
  "image",
  "title",
  "text",
 }
 writer.Write(headers)

 for _, data := range scrapData {
  // Converting a data to an array of strings
  record := []string{
   data.title,
   data.description,
   data.price,
  }
  writer.Write(record)
 }

 writer.Flush()

 if err := writer.Error(); err != nil {
  log.Fatalln("Error writing CSV:", err)
 }

 return scrapData
}

func main() {
    scrapeAndWriteCSV()
}
